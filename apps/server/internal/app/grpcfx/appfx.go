package grpcfx

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	authMd "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/grpc/interceptor"
	"github.com/taskemapp/server/apps/server/internal/pkg/logger"
	"github.com/taskemapp/server/apps/server/internal/repository/token"
	v1 "github.com/taskemapp/server/apps/server/tools/gen/grpc/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"strconv"
)

const module = "grpc"

var App = fx.Options(
	fx.Module(
		module,
		fx.Decorate(
			func(l logger.Logger) logger.Logger {
				return l.WithScope(module)
			},
		),

		fx.Provide(
			fx.Private,
			interceptor.New,
			fx.Annotate(token.NewClient, fx.As(new(token.Repository))),

			New,
		),

		fx.Invoke(
			func(lc fx.Lifecycle, log *zap.Logger, c config.Config, server GrpcServer) {
				lc.Append(
					fx.Hook{
						OnStart: func(ctx context.Context) error {
							go func() {
								err := server.Run(c)
								log.Error("server stopped", zap.Error(err))
							}()

							log.Info("server started on port", zap.String("addr", strconv.Itoa(c.GrpcPort)))

							return nil
						},
						OnStop: func(ctx context.Context) error {
							log.Info("gracefully stopping grpc server")
							server.GracefulStop()

							log.Info("server stopped")
							return nil
						},
					},
				)
			},
		),
	),
)

type Opts struct {
	fx.In
	AuthServer    v1.AuthServer
	ProfileServer v1.ProfileServer
	TeamServer    v1.TeamServer
	Log           logger.Logger
	Ic            *interceptor.Interceptor
}

type GrpcServer struct {
	Srv *grpc.Server
}

func New(opts Opts) GrpcServer {
	logOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall, logging.FinishCall,
			logging.PayloadSent, logging.PayloadReceived,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			opts.Log.Error("Recovered from panic", zap.Any("panic", p))
			return status.Error(codes.Internal, "Internal server error")
		}),
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			selector.UnaryServerInterceptor(authMd.UnaryServerInterceptor(opts.Ic.Auth), selector.MatchFunc(opts.Ic.AuthMatcher)),
			opts.Ic.ProvideRID(),
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(interceptorLogger(opts.Log), logOpts...),
		),
		grpc.ChainStreamInterceptor(
			selector.StreamServerInterceptor(authMd.StreamServerInterceptor(opts.Ic.Auth), selector.MatchFunc(opts.Ic.AuthMatcher)),
			recovery.StreamServerInterceptor(recoveryOpts...),
			logging.StreamServerInterceptor(interceptorLogger(opts.Log), logOpts...),
		),
	)

	v1.RegisterAuthServer(srv, opts.AuthServer)
	v1.RegisterProfileServer(srv, opts.ProfileServer)
	v1.RegisterTeamServer(srv, opts.TeamServer)
	reflection.Register(srv)

	return GrpcServer{Srv: srv}
}

func (a GrpcServer) Run(c config.Config) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", c.GrpcPort))
	if err != nil {
		return errors.Wrap(err, "run")
	}
	if err = a.Srv.Serve(l); err != nil {
		return errors.Wrap(err, "run")
	}

	return nil
}

func (a GrpcServer) GracefulStop() {
	a.Srv.GracefulStop()
}

// interceptorLogger Retrieved from
// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/62b7de50cda5a5d633f1013bfbe50e0f38db34ef/interceptors/logging/examples/zap/example_test.go#L17
func interceptorLogger(l logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		log := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			log.Debug(msg)
		case logging.LevelInfo:
			log.Info(msg)
		case logging.LevelWarn:
			log.Warn(msg)
		case logging.LevelError:
			log.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

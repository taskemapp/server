package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	authserver "taskem-server/internal/grpc/auth"
	v1 "taskem-server/tools/gen/grpc/v1"
)

type Opts struct {
	fx.In
	AuthServer *authserver.Server
	Log        *zap.Logger
}

type App struct {
	fx.Out
	Srv *grpc.Server
}

func New(opts Opts) App {
	logOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall, logging.FinishCall,
			logging.PayloadSent, logging.PayloadReceived,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			opts.Log.Sugar().Errorw("Recovered from panic", "panic", p)
			return status.Error(codes.Internal, "Internal server error")
		}),
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(interceptorLogger(opts.Log), logOpts...),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recoveryOpts...),
			logging.StreamServerInterceptor(interceptorLogger(opts.Log), logOpts...),
		),
	)

	v1.RegisterAuthServer(srv, opts.AuthServer)

	return App{Srv: srv}
}

// interceptorLogger Retrieved from
// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/62b7de50cda5a5d633f1013bfbe50e0f38db34ef/interceptors/logging/examples/zap/example_test.go#L17
func interceptorLogger(l *zap.Logger) logging.Logger {
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

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

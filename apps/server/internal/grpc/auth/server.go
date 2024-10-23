package auth

import (
	"context"
	"errors"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/grpc/interceptor"
	"github.com/taskemapp/server/apps/server/internal/pkg/validation"
	"github.com/taskemapp/server/apps/server/internal/repository/token"
	"github.com/taskemapp/server/apps/server/internal/repository/user"
	"github.com/taskemapp/server/apps/server/internal/service/auth"
	"github.com/taskemapp/server/apps/server/tools/gen/grpc/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"regexp"
)

type Opts struct {
	fx.In
	Auth      auth.Service
	Logger    *zap.Logger
	Config    config.Config
	TokenRepo token.Repository
}

type Server struct {
	v1.UnimplementedAuthServer
	auth      auth.Service
	logger    *zap.Logger
	config    config.Config
	tokenRepo token.Repository
}

func New(opts Opts) *Server {
	return &Server{
		auth:      opts.Auth,
		logger:    opts.Logger,
		config:    opts.Config,
		tokenRepo: opts.TokenRepo,
	}
}

func (s *Server) Login(
	ctx context.Context,
	req *v1.LoginRequest,
) (*v1.LoginResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing argument: email")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing argument: password")
	}

	resp, err := s.auth.Login(
		ctx,
		auth.LoginOpts{
			Email:    req.Email,
			Password: req.Password,
		})

	if err != nil {
		s.logger.Sugar().Error(err)
		switch {
		case errors.Is(err, user.ErrNotFound):
			return nil, status.Error(codes.NotFound, "Not found")
		case errors.Is(err, auth.ErrTokenGen):
			return nil, status.Error(codes.Internal, "Token generate failed")
		case errors.Is(err, auth.ErrPwdMatch):
			return nil, status.Error(codes.InvalidArgument, "Wrong password")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &v1.LoginResponse{
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		TokenType:    resp.TokenType,
	}, nil
}

func (s *Server) SignUp(
	ctx context.Context,
	req *v1.SignupRequest,
) (*emptypb.Empty, error) {
	if req.UserName == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing argument: user_name")
	}

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing argument: email")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing argument: password")
	}

	isValidMail, err := regexp.MatchString(validation.EmailRegex, req.Email)
	if !isValidMail || err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid email: use format example@example.com")
	}

	if !validation.IsPwdComplex(req.Password) {
		return nil, status.Error(codes.InvalidArgument, "Password is too weak")
	}

	err = s.auth.Registration(
		ctx,
		auth.RegistrationOpts{
			Email:    req.Email,
			Name:     req.UserName,
			Password: req.Password,
		})
	if err != nil {
		s.logger.Sugar().Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) RefreshToken(
	ctx context.Context,
	req *v1.RefreshTokenRequest,
) (*v1.RefreshTokenResponse, error) {

	uid, err := interceptor.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := s.auth.RefreshToken(
		ctx,
		auth.RefreshTokenOpts{
			UserID: uid,
			Token:  req.Token,
		},
	)

	if err != nil {
		s.logger.Error("refresh token failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &v1.RefreshTokenResponse{
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		TokenType:    resp.TokenType,
	}, nil
}

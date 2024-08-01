package auth

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"regexp"
	"taskem-server/internal/pkg/validation"
	"taskem-server/internal/repositories/user"
	"taskem-server/internal/service/auth"
	"taskem-server/tools/gen/grpc/v1"
)

type Opts struct {
	fx.In
	Auth   auth.Service
	Logger *zap.Logger
}

type Server struct {
	v1.UnimplementedAuthServer
	auth   auth.Service
	logger *zap.Logger
}

func New(opts Opts) *Server {
	return &Server{
		auth:   opts.Auth,
		logger: opts.Logger,
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
		return nil, err
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
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

package auth

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"taskem/internal/repositories/user"
	"taskem/internal/service/auth"
	authv1 "taskem/tools/gen/grpc/v1/auth"
)

type Opts struct {
	fx.In
	Auth auth.Auth
}

type Server struct {
	authv1.UnimplementedAuthServer
	Auth auth.Auth
}

func New(opts Opts) *Server {
	return &Server{Auth: opts.Auth}
}

func (s *Server) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {

	resp, err := s.Auth.Login(
		ctx,
		auth.LoginOpts{
			Email:    req.Email,
			Password: req.Password,
		})

	if err != nil {
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

	return &authv1.LoginResponse{
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (s *Server) SignUp(
	ctx context.Context,
	req *authv1.SignupRequest,
) (*emptypb.Empty, error) {

	err := s.Auth.Registration(
		ctx,
		auth.RegistrationOpts{
			Email:    req.Email,
			Name:     req.UserName,
			Password: req.Password,
		})
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &emptypb.Empty{}, nil
}

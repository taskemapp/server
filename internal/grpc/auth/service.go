package auth

import (
	"context"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	authv1 "taskem/tools/gen/grpc/v1/auth"
)

type authServer struct {
	fx.Out
	authv1.UnimplementedAuthServer
	auth Auth
}

func (s *authServer) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	if !strings.Contains(req.Email, "@") {
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	}

	resp, err := s.auth.Login(
		LoginOpts{
			email:    req.Email,
			password: req.Password,
		})
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{UserName: "ripls", Message: "asdasda", SessionId: "asdsada"}, nil
}

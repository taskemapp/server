package auth

import "context"

type Auth interface {
	Login(ctx context.Context, opts LoginOpts) (resp *LoginResponse, err error)
	Registration(ctx context.Context, opts RegistrationOpts) error
}

type LoginOpts struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token        string
	RefreshToken string
	TokenType    string
}

type RegistrationOpts struct {
	Email    string
	Name     string
	Password string
}

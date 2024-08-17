package auth

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	Login(ctx context.Context, opts LoginOpts) (resp *LoginResponse, err error)
	Registration(ctx context.Context, opts RegistrationOpts) error
	RefreshToken(ctx context.Context, opts RefreshTokenOpts) (resp *LoginResponse, err error)
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

type RefreshTokenOpts struct {
	UserID uuid.UUID
}

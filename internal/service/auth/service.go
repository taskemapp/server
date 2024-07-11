package auth

import (
	"context"
	"go.uber.org/fx"
	"taskem/internal/repositories/user"
)

type Opts struct {
	fx.In
	UserRepo user.Repository
}

type Service struct {
	UserRepo user.Repository
}

func New(opts Opts) *Service {
	return &Service{UserRepo: opts.UserRepo}
}

func (s *Service) Login(ctx context.Context, opts LoginOpts) (resp *LoginResponse, err error) {
	_, err = s.UserRepo.FindByEmail(ctx, opts.Email)
	if err != nil {
		return nil, err
	}

	//TODO impl passwd matching

	//TODO impl token generating
	return &LoginResponse{
		Token:        "mock",
		RefreshToken: "mock",
	}, nil
}

func (s *Service) Registration(ctx context.Context, opts RegistrationOpts) error {

	//TODO impl argon2 hashing
	passwd := "hashing with argon2"

	_, err := s.UserRepo.Create(ctx, user.CreateOpts{
		Email:    opts.Email,
		Name:     opts.Name,
		Password: passwd,
	})
	if err != nil {
		return err
	}

	return nil
}

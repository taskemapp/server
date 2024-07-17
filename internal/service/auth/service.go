package auth

import (
	"context"
	"github.com/alexedwards/argon2id"
	"go.uber.org/fx"
	"taskem/internal/config"
	"taskem/internal/pkg/jwt"
	"taskem/internal/repositories/user"
)

type Opts struct {
	fx.In
	UserRepo user.Repository
	Config   config.Config
}

type Service struct {
	UserRepo user.Repository
	Config   config.Config
}

func New(opts Opts) *Service {
	return &Service{
		UserRepo: opts.UserRepo,
		Config:   opts.Config,
	}
}

func (s *Service) Login(ctx context.Context, opts LoginOpts) (resp *LoginResponse, err error) {
	u, err := s.UserRepo.FindByEmail(ctx, opts.Email)
	if err != nil {
		return nil, err
	}

	match, err := argon2id.ComparePasswordAndHash(opts.Password, u.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, ErrPwdMatch
	}

	token, err := jwt.NewToken(jwt.Opts{
		Id:       u.ID,
		Duration: s.Config.TokenTtl,
		Email:    u.Email,
		Secret:   s.Config.TokenSecret,
	})
	if err != nil {
		return nil, err
	}

	refresh, err := jwt.NewToken(jwt.Opts{
		Id:       u.ID,
		Duration: s.Config.RefreshTokenTtl,
		Secret:   s.Config.TokenSecret,
	})

	if err != nil {
		return nil, ErrTokenGen
	}

	return &LoginResponse{
		Token:        token,
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}, nil
}

func (s *Service) Registration(ctx context.Context, opts RegistrationOpts) error {
	passwd, err := argon2id.CreateHash(opts.Password, argon2id.DefaultParams)
	if err != nil {
		return ErrPwdHash
	}

	_, err = s.UserRepo.Create(ctx, user.CreateOpts{
		Email:       opts.Email,
		Name:        opts.Name,
		Password:    passwd,
		DisplayName: opts.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

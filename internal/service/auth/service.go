package auth

import (
	"context"
	"github.com/alexedwards/argon2id"
	"go.uber.org/fx"
	"taskem-server/internal/config"
	"taskem-server/internal/pkg/jwt"
	"taskem-server/internal/repositories/user"
)

type Opts struct {
	fx.In
	UserRepo user.Repository
	Config   config.Config
}

type Auth struct {
	userRepo user.Repository
	config   config.Config
}

func New(opts Opts) *Auth {
	return &Auth{
		userRepo: opts.UserRepo,
		config:   opts.Config,
	}
}

func (a *Auth) Login(ctx context.Context, opts LoginOpts) (resp *LoginResponse, err error) {
	u, err := a.userRepo.FindByEmail(ctx, opts.Email)
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
		ID:       u.ID,
		Duration: a.config.TokenTtl,
		Email:    u.Email,
		Secret:   a.config.TokenSecret,
	})
	if err != nil {
		return nil, err
	}

	refresh, err := jwt.NewToken(jwt.Opts{
		ID:       u.ID,
		Duration: a.config.RefreshTokenTtl,
		Secret:   a.config.TokenSecret,
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

func (a *Auth) Registration(ctx context.Context, opts RegistrationOpts) error {
	passwd, err := argon2id.CreateHash(opts.Password, argon2id.DefaultParams)
	if err != nil {
		return ErrPwdHash
	}

	_, err = a.userRepo.Create(ctx, user.CreateOpts{
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

func (a *Auth) RefreshToken(ctx context.Context, opts RefreshTokenOpts) (resp *LoginResponse, err error) {
	u, err := a.userRepo.FindByID(ctx, opts.UserID)

	token, err := jwt.NewToken(jwt.Opts{
		ID:       u.ID,
		Duration: a.config.TokenTtl,
		Email:    u.Email,
		Secret:   a.config.TokenSecret,
	})
	if err != nil {
		return nil, err
	}

	refresh, err := jwt.NewToken(jwt.Opts{
		ID:       u.ID,
		Duration: a.config.RefreshTokenTtl,
		Secret:   a.config.TokenSecret,
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

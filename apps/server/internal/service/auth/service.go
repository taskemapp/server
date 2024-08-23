package auth

import (
	"context"
	"encoding/json"
	"github.com/alexedwards/argon2id"
	"github.com/taskemapp/server/apps/notification/pkg/notifier"
	"github.com/taskemapp/server/apps/server/internal/broker"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/pkg/jwt"
	"github.com/taskemapp/server/apps/server/internal/repositories/user"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	UserRepo user.Repository
	Config   config.Config
	Mq       broker.Mq
}

type Auth struct {
	userRepo user.Repository
	config   config.Config
	mq       broker.Mq
}

func New(opts Opts) *Auth {
	return &Auth{
		userRepo: opts.UserRepo,
		config:   opts.Config,
		mq:       opts.Mq,
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

	body, err := json.Marshal(notifier.Notification{
		Type:      notifier.EmailNotify,
		Recipient: opts.Email,
		Subject:   "OTP",
		Message:   "otp",
	})
	if err != nil {
		return nil, err
	}
	err = a.mq.Send(broker.SendOpts{
		Body: body,
	}, broker.NotificationChannel)

	if err != nil {
		return nil, err
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

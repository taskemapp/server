package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/taskemapp/server/apps/notification/pkg/notifier"
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/pkg/jwt"
	"github.com/taskemapp/server/apps/server/internal/repositories/token"
	"github.com/taskemapp/server/apps/server/internal/repositories/user"
	"github.com/taskemapp/server/libs/queue"
	"github.com/taskemapp/server/libs/template"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	TokenRepo token.Repository
	UserRepo  user.Repository
	Config    config.Config
	Br        queue.Queue
}

type Auth struct {
	userRepo  user.Repository
	config    config.Config
	br        queue.Queue
	tokenRepo token.Repository
}

func New(opts Opts) *Auth {
	return &Auth{
		br:        opts.Br,
		userRepo:  opts.UserRepo,
		config:    opts.Config,
		tokenRepo: opts.TokenRepo,
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

	access, err := jwt.NewToken(jwt.Opts{
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

	err = a.tokenRepo.SetToken(ctx, token.CreateOpts{
		ID:        u.ID,
		TokenType: "access",
		Token:     access,
		Duration:  a.config.TokenTtl,
	})
	if err != nil {
		return nil, err
	}

	err = a.tokenRepo.SetToken(ctx, token.CreateOpts{
		ID:        u.ID,
		TokenType: "refresh",
		Token:     refresh,
		Duration:  a.config.RefreshTokenTtl,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        access,
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

	temp, err := template.Get(template.VerifyEmailTemplate)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	err = temp.Execute(&buff, template.VerifyEmail{
		Name:             opts.Name,
		ConfirmationLink: "conf-link",
		UnsubscribeLink:  "unsubscribe-link",
	})
	if err != nil {
		return err
	}

	body, err := json.Marshal(notifier.EmailNotification{
		Notification: notifier.Notification{
			Title:   "Verify email",
			Message: buff.String(),
		},
		To:   opts.Email,
		From: "no-replay@taskem.test",
	})
	if err != nil {
		return err
	}

	err = a.br.Publish(notifier.ChannelEmail, queue.Message{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *Auth) RefreshToken(ctx context.Context, opts RefreshTokenOpts) (resp *LoginResponse, err error) {
	_, err = a.tokenRepo.GetToken(ctx, fmt.Sprintf("%s:%s", jwt.Refresh, opts.UserID))
	if err != nil {

		return nil, err
	}

	u, err := a.userRepo.FindByID(ctx, opts.UserID)
	if err != nil {
		return nil, err
	}

	access, err := jwt.NewToken(jwt.Opts{
		ID:       u.ID,
		Duration: a.config.TokenTtl,
		Email:    u.Email,
		Secret:   a.config.TokenSecret,
	})
	if err != nil {
		return nil, ErrTokenGen
	}

	refresh, err := jwt.NewToken(jwt.Opts{
		ID:       u.ID,
		Duration: a.config.RefreshTokenTtl,
		Secret:   a.config.TokenSecret,
	})

	if err != nil {
		return nil, ErrTokenGen
	}

	err = a.tokenRepo.SetToken(ctx, token.CreateOpts{
		ID:        u.ID,
		TokenType: "access",
		Token:     access,
		Duration:  a.config.TokenTtl,
	})
	if err != nil {
		return nil, err
	}

	err = a.tokenRepo.SetToken(ctx, token.CreateOpts{
		ID:        u.ID,
		TokenType: "refresh",
		Token:     refresh,
		Duration:  a.config.RefreshTokenTtl,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
	}, nil
}

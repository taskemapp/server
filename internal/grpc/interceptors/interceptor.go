package interceptors

import (
	"go.uber.org/fx"
	"taskem-server/internal/config"
	"taskem-server/internal/repositories/token"
)

type Opts struct {
	fx.In
	Config    config.Config
	TokenRepo token.Repository
}

type Interceptor struct {
	c         config.Config
	tokenRepo token.Repository
}

func New(opts Opts) *Interceptor {
	return &Interceptor{
		c:         opts.Config,
		tokenRepo: opts.TokenRepo,
	}
}

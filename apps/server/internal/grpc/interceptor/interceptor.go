package interceptor

import (
	"github.com/taskemapp/server/apps/server/internal/config"
	"github.com/taskemapp/server/apps/server/internal/pkg/logger"
	"github.com/taskemapp/server/apps/server/internal/repository/token"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	Config    config.Config
	TokenRepo token.Repository
	Logger    logger.Logger
}

type Interceptor struct {
	c         config.Config
	tokenRepo token.Repository
	logger    logger.Logger
}

func New(opts Opts) *Interceptor {
	return &Interceptor{
		c:         opts.Config,
		tokenRepo: opts.TokenRepo,
		logger:    opts.Logger,
	}
}

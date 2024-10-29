package token

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/go-redis/redis/v8"
	"github.com/taskemapp/server/apps/server/internal/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Opts struct {
	fx.In
	Client *redis.Client
	Logger logger.Logger
}

type Client struct {
	client *redis.Client
	logger logger.Logger
}

func NewClient(opts Opts) (*Client, error) {
	return &Client{
		client: opts.Client,
		logger: opts.Logger,
	}, nil
}

func (rc *Client) SetToken(ctx context.Context, opts CreateOpts) error {
	val, err := rc.client.Set(
		ctx,
		fmt.Sprintf("%s:%s", opts.TokenType, opts.ID.String()),
		opts.Token,
		opts.Duration,
	).Result()

	if err != nil {
		rc.logger.Error("Failed to set token: ", zap.Error(err))
		return err
	}

	rc.logger.Info("Token set: %s", zap.String("ok", val))
	return nil
}

func (rc *Client) GetToken(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		rc.logger.Warn("Failed to get token: ", zap.Error(err))
		return "", errors.Wrap(ErrNotFound, "Failed to get token")
	}
	if err != nil {
		rc.logger.Error("Failed to get token: ", zap.Error(err))
		return "", errors.Wrap(err, "Failed to get token")
	}

	return val, err
}

func (rc *Client) DelToken(ctx context.Context, key string) error {
	return errors.Wrap(rc.client.Del(ctx, key).Err(), "Failed to delete token")
}

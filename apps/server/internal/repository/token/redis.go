package token

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Opts struct {
	fx.In
	Client *redis.Client
	Logger *zap.Logger
}

type Client struct {
	client *redis.Client
	logger *zap.Logger
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
		rc.logger.Sugar().Error("Failed to set token: ", err)
		return err
	}

	rc.logger.Sugar().Infof("Token set: %s", val)
	return nil
}

func (rc *Client) GetToken(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		rc.logger.Sugar().Warn("Failed to get token: ", err)
		return "", errors.Wrap(ErrNotFound, "Failed to get token")
	}
	if err != nil {
		rc.logger.Sugar().Error("Failed to get token: ", err)
		return "", errors.Wrap(err, "Failed to get token")
	}

	return val, err
}

func (rc *Client) DelToken(ctx context.Context, key string) error {
	return errors.Wrap(rc.client.Del(ctx, key).Err(), "Failed to delete token")
}

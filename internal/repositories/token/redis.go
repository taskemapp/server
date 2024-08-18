package token

import (
	"context"
	"errors"
	"fmt"
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
	return rc.client.Set(ctx, fmt.Sprintf("%s:%s", opts.TokenType, opts.ID.String()), opts.Token, opts.Duration).Err()
}

func (rc *Client) GetToken(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	return val, err
}

func (rc *Client) DelToken(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

package token

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	GetToken(ctx context.Context, key string) (string, error)
	SetToken(ctx context.Context, opts CreateOpts) error
	DelToken(ctx context.Context, key string) error
}

type CreateOpts struct {
	ID        uuid.UUID
	TokenType string
	Token     string
	Duration  time.Duration
}

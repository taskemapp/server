package token

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	// GetToken get auth token from `in memory storage`, where key defined as
	// tokenType:userID
	GetToken(ctx context.Context, key string) (string, error)
	// SetToken store auth token in `in memory storage`, where key defined as
	// tokenType:userID
	SetToken(ctx context.Context, opts CreateOpts) error
	DelToken(ctx context.Context, key string) error
}

type CreateOpts struct {
	ID        uuid.UUID
	TokenType string
	Token     string
	Duration  time.Duration
}

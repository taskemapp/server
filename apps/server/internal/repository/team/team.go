package team

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, opts CreateOpts) (*Team, error)
	FindByID(ctx context.Context, taskID uuid.UUID) (*Team, error)
	FindMany(ctx context.Context, opts FindManyOpts) (*FindManyResult, error)
	FindManyBelongToUser(ctx context.Context, userID uuid.UUID, opts FindManyOpts) (*FindManyResult, error)
	Update(ctx context.Context, taskID uuid.UUID, opts UpdateOpts) (*Team, error)
}

type CreateOpts struct {
	Name        string
	Description string
	Creator     uuid.UUID
}

type UpdateOpts struct {
	Name        *string
	Description *string
	HeaderImage *string
	Creator     *uuid.UUID
}

type FindManyOpts struct {
	Page    int
	PerPage int
}

type FindManyResult struct {
	Teams       []Team
	Total       int
	HasNextPage bool
}

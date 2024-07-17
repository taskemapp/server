package user

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, opts CreateOpts) (*User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindManyBelongToTeam(ctx context.Context, teamID uuid.UUID, opts FindManyOpts) (*FindManyResult, error)
	FindMany(ctx context.Context, opts FindManyOpts) (*FindManyResult, error)
	Update(ctx context.Context, userID uuid.UUID, opts UpdateOpts) (*User, error)
	DeleteAccount(ctx context.Context, userID uuid.UUID) error
}

type FindManyOpts struct {
	Page    int
	PerPage int
}

type FindManyResult struct {
	Users       []User
	Total       int
	HasNextPage bool
}

type CreateOpts struct {
	Name        string
	DisplayName string
	Email       string
	Password    string
}

type UpdateOpts struct {
	DisplayName *string
	Email       *string
	AvatarUrl   *string
	IsVerified  *string
}

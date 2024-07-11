package user

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, opts CreateOpts) (user *User, err error)
	FindByID(ctx context.Context, userID uuid.UUID) (user *User, err error)
	FindByName(ctx context.Context, name string) (user *User, err error)
	FindByEmail(ctx context.Context, email string) (user *User, err error)
	Update(ctx context.Context, userID uuid.UUID, opts UpdateOpts) (user *User, err error)
	DeleteAccount(ctx context.Context, userID uuid.UUID) error
	FindMany(ctx context.Context, opts FindManyOpts) (res *FindManyResult, err error)
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
	AvatarUrl   *string
}

type UpdateOpts struct {
	DisplayName *string
	Email       *string
	AvatarUrl   *string
	IsVerified  *string
}

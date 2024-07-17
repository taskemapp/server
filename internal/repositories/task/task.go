package task

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	Create(ctx context.Context, opts CreateOpts) (*Task, error)
	FindByID(ctx context.Context, taskID uuid.UUID) (*Task, error)
	FindMany(ctx context.Context, opts FindManyOpts) (*FindManyResult, error)
	FindManyBelongToTeam(ctx context.Context, teamID uuid.UUID, opts FindManyOpts) (*FindManyResult, error)
	FindManyBelongToUser(ctx context.Context, userID uuid.UUID, opts FindManyOpts) ([]Task, error)
	Update(ctx context.Context, taskID uuid.UUID, opts UpdateOpts) (*Task, error)
}

type CreateOpts struct {
	Name           string
	Description    string
	Status         string
	TeamID         uuid.UUID
	AssignedUserID *uuid.UUID
	Creator        uuid.UUID
	EndAt          *time.Time
}

type UpdateOpts struct {
	Name           *string
	Description    *string
	Status         *string
	AssignedUserID *uuid.UUID
	EndAt          *time.Time
}

type FindManyOpts struct {
	Status string
}

type FindManyResult struct {
	Tasks       []Task
	Total       int
	HasNextPage bool
}

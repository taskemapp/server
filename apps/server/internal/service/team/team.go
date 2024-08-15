package team

import (
	"context"
	"github.com/google/uuid"
	"server/internal/repositories/team"
	"server/internal/service"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (*team.Team, error)
	GetUserTeams(ctx context.Context, userID uuid.UUID, pgOpts service.PaginationOpts) (*team.FindManyResult, error)
	Create(ctx context.Context, opts CreateOpts) (*team.Team, error)
	Join(ctx context.Context, opts JoinOpts) error
	Leave(ctx context.Context, opts LeaveOpts) error
	// GetRoles TODO(ripls56): change returning signature
	//
	// Deprecated: not recommended to use
	GetRoles(ctx context.Context, teamID uuid.UUID) error
}

type CreateOpts struct {
	CreatorID   uuid.UUID
	Name        string
	Description string
}

type JoinOpts struct {
	UserID uuid.UUID
	TeamID uuid.UUID
}

type LeaveOpts struct {
	UserID uuid.UUID
	TeamID uuid.UUID
}

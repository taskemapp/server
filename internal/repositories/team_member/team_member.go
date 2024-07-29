package team_member

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, tmID uuid.UUID) (*TeamMember, error)
	Create(ctx context.Context, opts CreateOpts) (*TeamMember, error)
	// Update row by team member pk id
	Update(ctx context.Context, tmID uuid.UUID, opts UpdateOpts) (*TeamMember, error)
}

type CreateOpts struct {
	UserID uuid.UUID
	TeamID uuid.UUID
}

type UpdateOpts struct {
	IsLeaved *bool
}

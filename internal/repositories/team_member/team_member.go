package team_member

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	Create(ctx context.Context, opts CreateOpts) (*TeamMember, error)
	// Update row by team member pk id
	Update(ctx context.Context, tmID uuid.UUID, opts UpdateOpts) (*TeamMember, error)
}

type CreateOpts struct {
	UserID uuid.UUID
	TeamID uuid.UUID
}

type UpdateOpts struct {
	LeaveAt *time.Time
	Leave   *bool
}

package team_member

import (
	"github.com/google/uuid"
	"time"
)

type TeamMember struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	TeamID   uuid.UUID
	JoinedAt time.Time
	LeavedAt *time.Time
	IsLeaved bool
}

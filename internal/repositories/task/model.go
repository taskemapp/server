package task

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID             uuid.UUID
	Name           string
	Description    string
	Status         string
	TeamID         uuid.UUID
	AssignedUserID *uuid.UUID
	Creator        uuid.UUID
	CreatedAt      time.Time
	EditedAt       *time.Time
	EndedAt        *time.Time
}

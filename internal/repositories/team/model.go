package team

import (
	"github.com/google/uuid"
	"time"
)

type Team struct {
	ID             uuid.UUID
	Name           string
	Description    string
	HeaderImageUrl *string
	Creator        uuid.UUID
	CreatedAt      time.Time
	EditedAt       *time.Time
}

package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Email       string
	AvatarUrl   string
	CreatedAt   time.Time
	EditedAt    time.Time
}

package user_files

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, opts CreateUserFileOpts) (*UserFile, error)
	Update(ctx context.Context, opts UpdateUserFileOpts) (*UserFile, error)
}

type CreateUserFileOpts struct {
	UserID   uuid.UUID
	CdnPath  string
	FileName string
	MimeType string
}

type UpdateUserFileOpts struct {
	UserID   uuid.UUID
	CdnPath  string
	FileName string
	MimeType string
}

type UserFile struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	CdnPath  string
	FileName string
	MimeType string
}

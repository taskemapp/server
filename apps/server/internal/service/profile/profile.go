package profile

import (
	"context"
	"github.com/google/uuid"
)

type Service interface {
	UploadAvatar(ctx context.Context, userID uuid.UUID, opts UploadAvatarOpts) error
}

type UploadAvatarOpts struct {
	Avatar []byte
}
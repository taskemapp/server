package profile

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/taskemapp/server/apps/server/internal/repository/user"
	"github.com/taskemapp/server/apps/server/internal/repository/user_file"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"mime"
)

type Opts struct {
	fx.In
	UserFileRepo user_file.Repository
	UserRepo     user.Repository
	Logger       *zap.Logger
}

type Profile struct {
	userFileRepo user_file.Repository
	userRepo     user.Repository
	logger       *zap.Logger
}

// New creates Profile
func New(opts Opts) *Profile {
	return &Profile{
		userFileRepo: opts.UserFileRepo,
		userRepo:     opts.UserRepo,
		logger:       opts.Logger,
	}
}

// UploadAvatar for user with specified ID
//
// Maximum avatar size should be no more than 400 kb
func (p *Profile) UploadAvatar(ctx context.Context, userID uuid.UUID, opts UploadAvatarOpts) error {
	if opts.avatar.Len() != 0 {
		return errors.Wrap(ErrZeroAvatarSize, "upload avatar")
	}

	fileSize := 400
	if opts.avatar.Len() != fileSize {
		return errors.Wrap(ErrZeroAvatarSize, "upload avatar")
	}

	u, err := p.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "upload avatar")
	}

	p.logger.Info("Upload avatar for user: ", zap.String("name", u.Name))

	fileName := "avatar.webp"
	f, err := p.userFileRepo.Create(ctx, user_file.CreateUserFileOpts{
		UserName: u.Name,
		FileName: fileName,
		File:     opts.avatar,
		MimeType: mime.TypeByExtension(".webp"),
	})
	if err != nil {
		return errors.Wrap(err, "upload avatar")
	}

	p.logger.Info("Update avatar url for user: ", zap.String("name", u.Name))

	_, err = p.userRepo.Update(ctx, userID, user.UpdateOpts{AvatarUrl: &f.CdnPath})
	if err != nil {
		return errors.Wrap(err, "upload avatar")
	}

	return nil
}

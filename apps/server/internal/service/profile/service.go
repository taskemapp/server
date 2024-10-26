package profile

import (
	"bytes"
	"context"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/taskemapp/server/apps/server/internal/repository/user"
	"github.com/taskemapp/server/apps/server/internal/repository/user_file"
	"github.com/taskemapp/server/apps/server/internal/service/profile/image"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"mime"
)

type Opts struct {
	fx.In
	UserFileRepo user_file.Repository
	UserRepo     user.Repository
	Logger       *zap.Logger
	Processing   image.Processing
}

type Profile struct {
	userFileRepo user_file.Repository
	userRepo     user.Repository
	logger       *zap.Logger
	processing   image.Processing
}

// New creates Profile
func New(opts Opts) *Profile {
	u := opts.UserRepo
	if u == nil {
		panic("user repo didnt provided")
	}
	return &Profile{
		userFileRepo: opts.UserFileRepo,
		userRepo:     opts.UserRepo,
		logger:       opts.Logger,
		processing:   opts.Processing,
	}
}

// UploadAvatar for user with specified ID
//
// Maximum avatar size should be no more than 1 mb
func (p *Profile) UploadAvatar(ctx context.Context, userID uuid.UUID, opts UploadAvatarOpts) error {
	avatar := opts.Avatar
	if len(avatar) == 0 {
		return errors.Wrap(ErrZeroAvatarSize, "upload avatar")
	}

	fileSize := 1024 * 1024
	if len(avatar) > fileSize {
		return errors.Wrap(ErrWrongAvatarSize, "upload avatar")
	}

	err := p.processing.ConvertToWebp(avatar)
	if err != nil {
		return errors.Wrap(err, "upload avatar")
	}
	p.logger.Info("image converted", zap.Int("len", len(avatar)))

	u, err := p.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "upload avatar")
	}
	p.logger.Info("user found", zap.Int("len", len(avatar)))

	fileName := "avatar.webp"
	var buff bytes.Buffer
	buff.Write(avatar)
	f, err := p.userFileRepo.Create(ctx, user_file.CreateUserFileOpts{
		UserName: u.Name,
		FileName: fileName,
		File:     buff,
		MimeType: mime.TypeByExtension(".webp"),
	})
	if err != nil {
		return errors.Wrap(err, "upload avatar")
	}
	p.logger.Info("avatar uploaded for user")

	_, err = p.userRepo.Update(ctx, userID, user.UpdateOpts{AvatarUrl: &f.CdnPath})
	if err != nil {
		return errors.Wrap(err, "upload avatar")
	}

	p.logger.Info("avatar url updated")

	return nil
}

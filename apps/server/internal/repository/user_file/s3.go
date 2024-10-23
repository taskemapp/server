package user_file

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/minio/minio-go/v7"
	"github.com/taskemapp/server/apps/server/internal/pkg/s3"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"strings"
)

var _ Repository = (*File)(nil)

type Opts struct {
	fx.In
	S3    *minio.Client
	S3Cfg s3.Config
}

type File struct {
	s3    *minio.Client
	s3Cfg s3.Config
}

func New(opts Opts) *File {
	return &File{
		s3:    opts.S3,
		s3Cfg: opts.S3Cfg,
	}
}

func (s File) Create(ctx context.Context, opts CreateUserFileOpts) (*UserFile, error) {
	file := opts.File
	filePath := strings.Join([]string{
		opts.UserName,
		opts.FileName,
	}, "/")

	_, err := s.s3.PutObject(
		ctx,
		s.s3Cfg.Bucket,
		filePath,
		&file,
		int64(file.Len()),
		minio.PutObjectOptions{
			ContentType:  opts.MimeType,
			AutoChecksum: minio.ChecksumSHA256,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "s3 create")
	}

	return &UserFile{
		CdnPath:  s.buildFileURL(opts.FileName),
		FileName: opts.FileName,
		MimeType: opts.MimeType,
	}, nil
}

func (s File) buildFileURL(filePath string) string {
	var sb strings.Builder
	var err error

	_, err = sb.WriteString(strings.TrimSuffix(s.s3Cfg.Host, "/"))
	err = multierr.Append(err, err)

	_, err = sb.WriteString("/")
	err = multierr.Append(err, err)

	_, err = sb.WriteString(s.s3Cfg.Bucket)
	err = multierr.Append(err, err)

	_, err = sb.WriteString("/")
	err = multierr.Append(err, err)

	_, err = sb.WriteString(filePath)
	err = multierr.Append(err, err)

	return sb.String()
}

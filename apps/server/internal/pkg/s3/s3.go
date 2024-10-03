package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New(cfg Config) (*minio.Client, error) {
	c, err := minio.New(cfg.Host, &minio.Options{
		Creds:  credentials.NewStaticV4("root", "password", ""),
		Secure: cfg.Secure,
	})
	if err != nil {
		return nil, errors.Wrap(err, "s3 new")
	}

	return c, nil
}

func Invoke(lc fx.Lifecycle, log *zap.Logger, cfg Config, c *minio.Client) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Info("Listing buckets")
				ctx, cancel := context.WithTimeout(ctx, time.Second*2)
				defer cancel()

				ok, err := c.BucketExists(ctx, cfg.Bucket)
				if err != nil {
					return fmt.Errorf("cannot list buckets: %w", err)
				}

				log.Sugar().Debug("Bucket exist: ", ok)

				if !ok {
					err = c.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
					if err != nil {
						return fmt.Errorf("cannot create bucket: %w", err)
					}
				}

				err = c.SetBucketPolicy(
					ctx,
					cfg.Bucket,
					`{
						"Version": "2012-10-17",
						"Statement": [
						{
							"Effect": "Allow",
							"Principal": "*",
							"Action": ["s3:GetObject"],
							"Resource": [
								"arn:aws:s3:::`+cfg.Bucket+`/team/*/header.webp",
								"arn:aws:s3:::`+cfg.Bucket+`/user/*/avatar.webp"
							]
						}
					]
				}`,
				)

				if err != nil {
					return fmt.Errorf("cannot set bucket policy: %w", err)
				}

				return nil
			},
			OnStop: nil,
		},
	)
}

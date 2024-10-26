package user_file

import (
	"bytes"
	"context"
)

type Repository interface {
	// Create loads the file into S3 storage and returns information about the file.
	Create(ctx context.Context, opts CreateUserFileOpts) (*UserFile, error)
}

type CreateUserFileOpts struct {
	UserName string
	FileName string
	File     bytes.Buffer
	// MimeType for set use mime package
	//
	// For example:
	// mime.TypeByExtension(".webp")
	MimeType string
}

type UserFile struct {
	CdnPath  string
	FileName string
	MimeType string
}

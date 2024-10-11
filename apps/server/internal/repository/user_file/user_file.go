package user_file

import (
	"bytes"
	"context"
)

type Repository interface {
	Create(ctx context.Context, opts CreateUserFileOpts) (*UserFile, error)
}

type CreateUserFileOpts struct {
	UserName string
	FilePath string
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

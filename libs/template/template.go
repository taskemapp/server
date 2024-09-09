package template

import (
	"github.com/go-faster/errors"
	"html/template"
	"os"
	"path/filepath"
)

type Template int

const (
	VerifyEmailTemplate Template = iota
)

const (
	verifyEmailName = "verify_email.gohtml"
)

type VerifyEmail struct {
	Name             string
	ConfirmationLink string
	UnsubscribeLink  string
}

func Get(t Template) (*template.Template, error) {
	wd, err := os.Getwd()
	path := filepath.Join(wd, "templates")
	if err != nil {
		return nil, err
	}

	switch t {
	case VerifyEmailTemplate:
		path = filepath.Join(path, verifyEmailName)
		file, err := template.ParseFiles(path)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse verify email template")
		}
		return file, nil
	default:
		return nil, errors.New("invalid template")
	}
}

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

type option struct {
	dir string
}

type OptionFn func(*option)

func WithDir(dir string) OptionFn {
	return func(o *option) {
		o.dir = dir
	}
}

func Get(t Template, opts ...OptionFn) (*template.Template, error) {
	wd, err := os.Getwd()
	defaultOpt := option{dir: wd}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&defaultOpt)
	}

	path := filepath.Join(defaultOpt.dir, "templates")
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

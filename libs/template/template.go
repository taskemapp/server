package template

import (
	"github.com/go-faster/errors"
	"html/template"
	"os"
	"path/filepath"
)

// Type represent template type
//
// Example:
// template.VerifyEmailTemplate
type Type int

const (
	VerifyEmailTemplate Type = iota + 1
)

const (
	verifyEmailName = "verify_email.gohtml"
)

type option struct {
	dir string
}

type OptionFn func(*option)

func WithDir(dir string) OptionFn {
	return func(o *option) {
		o.dir = dir
	}
}

type VerifyEmail struct {
	Name             string
	ConfirmationLink string
	UnsubscribeLink  string
}

// Get retrieves and parses the specified template.
// It accepts the TemplateType type and the variable number OptionFn to set the template's directory.
// If no directory is specified, the current working directory is used by default.
func Get(t Type, opts ...OptionFn) (*template.Template, error) {
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
			return nil, errors.Wrap(err, "verify email")
		}
		return file, nil
	default:
		return nil, errors.New("invalid template")
	}
}

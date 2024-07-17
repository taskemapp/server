package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

var (
	ErrMissingSecret = errors.New("missing arg: secret")
	ErrDuration      = errors.New("duration less then 1 min")
)

type Opts struct {
	Id       uuid.UUID
	Email    string
	Duration time.Duration
	Secret   string
}

// NewToken creates new JWT token for given user and app.
func NewToken(opts Opts) (token string, err error) {
	if opts.Secret == "" {
		return "", ErrMissingSecret
	}

	if opts.Duration < time.Minute {
		return "", ErrDuration
	}
	claims := jwt.MapClaims{
		"uid": opts.Id,
		"exp": time.Now().Add(opts.Duration).Unix(),
		"iat": time.Now().Unix(),
	}

	if opts.Email != "" {
		claims["email"] = opts.Email
	}

	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&claims,
	)

	tokenString, err := t.SignedString([]byte(opts.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

package jwt

import (
	"github.com/go-faster/errors"
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

func GetPayload(token string, secret string) (*jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Wrong token signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse jwt token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.Wrap(err, "Failed to get claims from token")
	}
	return &claims, nil
}

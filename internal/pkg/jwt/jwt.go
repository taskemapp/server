package jwt

import (
	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Opts struct {
	ID       uuid.UUID
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
		"uid": opts.ID,
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

// GetPayload undirected call `Validate` method,
// so if u use this method, u don't need to call `Validate` again
//
// Throws: ErrTokenValidation, ErrTokenParse
// if token expired throws ErrTokenExpired
func GetPayload(token string, secret string) (*jwt.MapClaims, error) {
	valid, err := Validate(token, secret)
	if err != nil {
		switch {
		case errors.Is(err, ErrTokenExpired):
			return nil, err
		case errors.Is(err, ErrTokenSigningValidation):
			return nil, ErrTokenValidation
		default:
			return nil, ErrTokenParse
		}
	}

	claims, ok := valid.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrTokenParse
	}
	return &claims, nil
}

// Validate provided token using secret
//
// Throws: ErrTokenParse, ErrTokenValidation, ErrTokenExpired
func Validate(token string, secret string) (*jwt.Token, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSigningValidation
		}
		return []byte(secret), nil
	})

	if parsed == nil {
		return nil, ErrTokenParse
	}

	switch {
	case parsed.Valid:
		return parsed, nil
	case !parsed.Valid:
		return nil, ErrTokenValidation
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, ErrTokenParse
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, ErrTokenExpired
	default:
		return nil, err
	}
}

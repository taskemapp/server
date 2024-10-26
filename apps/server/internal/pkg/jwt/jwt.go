package jwt

import (
	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// TokenType existing token types:
// Access, Refresh
//
// using for creating tokens with `type` in claims
type TokenType string

const (
	// Access token type see TokenType for more info
	Access TokenType = "access"
	// Refresh token type see TokenType for more info
	Refresh = "refresh"
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
		"uid":  opts.ID,
		"exp":  time.Now().Add(opts.Duration).Unix(),
		"iat":  time.Now().Unix(),
		"type": Refresh,
	}

	if opts.Email != "" {
		claims["email"] = opts.Email
		claims["type"] = Access
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

// Claims wrap jwt.MapClaims to get rid of import conflict
type Claims jwt.MapClaims

// GetPayload undirected call `Validate` method,
// so if u use this method, u don't need to call `Validate` again
//
// Throws: ErrTokenValidation, ErrTokenParse
// if token expired throws ErrTokenExpired
func GetPayload(token string, secret string) (*Claims, error) {
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

	c := Claims(claims)
	return &c, nil
}

// Token wrap jwt.Token to get rid of import conflict
type Token *jwt.Token

// Validate provided token using secret
//
// Throws: ErrTokenParse, ErrTokenValidation, ErrTokenExpired
func Validate(given string, secret string) (Token, error) {
	parsed, err := jwt.Parse(given, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSigningValidation
		}
		return []byte(secret), nil
	})

	if parsed == nil {
		return nil, ErrTokenParse
	}

	token := Token(parsed)

	switch {
	case token.Valid:
		return token, nil
	case !token.Valid:
		return nil, ErrTokenValidation
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, ErrTokenParse
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, ErrTokenExpired
	default:
		return nil, err
	}
}

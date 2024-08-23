package jwt

import "errors"

var (
	ErrMissingSecret          = errors.New("missing arg: secret")
	ErrDuration               = errors.New("duration less then 1 min")
	ErrTokenParse             = errors.New("failed to parse token")
	ErrTokenSigningValidation = errors.New("token signing validation error")
	ErrPayloadExtract         = errors.New("failed payload extract")
	ErrTokenValidation        = errors.New("failed token validation")
	ErrTokenExpired           = errors.New("token expired")
)

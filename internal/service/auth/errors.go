package auth

import "errors"

var (
	// ErrTokenGen failed generate token error
	ErrTokenGen = errors.New("token generation failed")
	// ErrPwdMatch wrong password error
	ErrPwdMatch = errors.New("password don't match")
	// ErrPwdHash failed password hashing
	ErrPwdHash = errors.New("failed password hashing")
)

package interceptor

import "errors"

var (
	ErrGetUserID = errors.New("failed to get uid from context")
	ErrRequestID = errors.New("failed to get rid from context")
)

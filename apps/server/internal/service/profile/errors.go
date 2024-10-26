package profile

import "errors"

var (
	ErrWrongAvatarSize = errors.New("avatar size is larger than the allowed size")
	ErrZeroAvatarSize  = errors.New("given avatar has zero size")
)

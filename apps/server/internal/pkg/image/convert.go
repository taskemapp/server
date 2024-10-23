package image

import (
	"github.com/go-faster/errors"
	"github.com/h2non/bimg"
)

func ConvertToWebp(b []byte) (i []byte, err error) {
	n := bimg.NewImage(b)
	t := n.Type()
	if t != "webp" {
		i, err = n.Convert(bimg.WEBP)
		if err != nil {
			return nil, errors.Wrap(err, "convert to webp")
		}
	}

	return i, nil
}

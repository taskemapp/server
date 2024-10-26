package image

import (
	"github.com/go-faster/errors"
	"github.com/h2non/bimg"
)

type ProfileProcessing struct {
}

func NewProcessing() *ProfileProcessing {
	return &ProfileProcessing{}
}

// TODO(ripls56): не работает с png, над зависимости глянуть в докере

func (p *ProfileProcessing) ConvertToWebp(b []byte) error {
	n := bimg.NewImage(b)
	t := n.Type()
	if t != "webp" {
		i, err := n.Convert(bimg.WEBP)
		if err != nil {
			return errors.Wrap(err, "convert to webp")
		}
		b = i
	}

	return nil
}

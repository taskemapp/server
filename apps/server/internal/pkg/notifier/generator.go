package notifier

import (
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"go.uber.org/multierr"
	"strings"
)

type LinkGenerator interface {
	VerifyLink() (link string, id uuid.UUID, err error)
	UnsubLink() (link string, id uuid.UUID, err error)
}

type BasicGenerator struct {
	HostDomain string
}

func (b *BasicGenerator) VerifyLink() (string, uuid.UUID, error) {
	var sb strings.Builder
	var err error

	confirmID, err := uuid.NewV7()
	if err != nil {
		return "", uuid.Nil, errors.Wrap(err, "confirm link")
	}

	_, err = sb.WriteString(b.HostDomain)
	err = multierr.Append(err, err)

	_, err = sb.WriteString("/verify?id=")
	err = multierr.Append(err, err)

	_, err = sb.WriteString(b.HostDomain)
	err = multierr.Append(err, err)

	return sb.String(), confirmID, err
}

func (b *BasicGenerator) UnsubLink() (string, uuid.UUID, error) {
	var sb strings.Builder
	var err error
	unsubID, err := uuid.NewV7()
	if err != nil {
		return "", uuid.Nil, errors.Wrap(err, "unsub link")
	}

	_, err = sb.WriteString(b.HostDomain)
	err = multierr.Append(err, err)

	_, err = sb.WriteString("/unsub?id=")
	err = multierr.Append(err, err)

	_, err = sb.WriteString(b.HostDomain)
	err = multierr.Append(err, err)

	return sb.String(), unsubID, err
}

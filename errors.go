package twitterscraper

import (
	"errors"
	"fmt"
)

type ErrRateLimitExceeded struct {
	description []byte
}

func (e ErrRateLimitExceeded) Error() string {
	return fmt.Sprintf("Rate limit exceeded: %s", e.description)
}

func (e ErrRateLimitExceeded) Is(err error) bool {
	return errors.As(err, &e)
}

type ErrOther struct {
	StatusCode  int
	status      string
	description []byte
}

func (o ErrOther) Error() string {
	return fmt.Sprintf("response status %s: %s", o.status, o.description)
}

func (o ErrOther) Is(err error) bool {
	return errors.As(err, &o)
}

var (
	ErrGuestTokenNotFound       = errors.New("guest_token not found")
	ErrDecodeGuestTokenResponse = errors.New("failed to decode guest token response")
)

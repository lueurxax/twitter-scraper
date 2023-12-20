package twitterscraper

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrRateLimitExceeded_Match(t *testing.T) {
	t.Run("is_match", func(t *testing.T) {
		e := ErrRateLimitExceeded{
			description: []byte("Rate limit exceeded: 123"),
		}

		assert.True(t, errors.Is(e, ErrRateLimitExceeded{}))
	})
	t.Run("is_not_match", func(t *testing.T) {
		e := ErrRateLimitExceeded{
			description: []byte("Rate limit exceeded: 123"),
		}

		assert.False(t, errors.Is(e, errors.New("some error")))
	})
}

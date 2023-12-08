package twitterscraper_test

import (
	"context"
	"testing"

	twitterscraper "github.com/lueurxax/twitter-scraper"
)

func TestGetGuestToken(t *testing.T) {
	scraper := twitterscraper.New()
	if err := scraper.GetGuestToken(context.Background()); err != nil {
		t.Errorf("getGuestToken() error = %v", err)
	}
	if !scraper.IsGuestToken() {
		t.Error("Expected non-empty guestToken")
	}
}

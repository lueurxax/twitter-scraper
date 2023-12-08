package twitterscraper_test

import (
	"context"
	"testing"
)

func TestGetTrends(t *testing.T) {
	if skipAuthTest {
		t.Skip("Skipping test due to environment variable")
	}

	ctx := context.Background()

	trends, err := testScraper.GetTrends(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(trends) != 20 {
		t.Errorf("Expected 20 trends, got %d: %#v", len(trends), trends)
	}

	for _, trend := range trends {
		if trend == "" {
			t.Error("Expected trend is empty")
		}
	}
}

package twitterscraper_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	twitterscraper "github.com/lueurxax/twitter-scraper"
)

var (
	username     = os.Getenv("TWITTER_USERNAME")
	password     = os.Getenv("TWITTER_PASSWORD")
	email        = os.Getenv("TWITTER_EMAIL")
	skipAuthTest = os.Getenv("SKIP_AUTH_TEST") != ""
	testScraper  = twitterscraper.New()
)

func init() {
	if username != "" && password != "" && !skipAuthTest {
		ctx := context.Background()
		err := testScraper.Login(ctx, username, password, email)
		if err != nil {
			panic(fmt.Sprintf("Login() error = %v", err))
		}
	}
}

func TestAuth(t *testing.T) {
	if skipAuthTest {
		t.Skip("Skipping test due to environment variable")
	}
	ctx := context.Background()
	scraper := twitterscraper.New()
	if err := scraper.Login(ctx, username, password, email); err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if !scraper.IsLoggedIn(ctx) {
		t.Fatalf("Expected IsLoggedIn() = true")
	}
	cookies := scraper.GetCookies()
	scraper2 := twitterscraper.New()
	scraper2.SetCookies(cookies)
	if !scraper2.IsLoggedIn(ctx) {
		t.Error("Expected restored IsLoggedIn() = true")
	}
	if err := scraper.Logout(ctx); err != nil {
		t.Errorf("Logout() error = %v", err)
	}
	if scraper.IsLoggedIn(ctx) {
		t.Error("Expected IsLoggedIn() = false")
	}
}

func TestLoginOpenAccount(t *testing.T) {
	scraper := twitterscraper.New()
	if err := scraper.LoginOpenAccount(context.Background()); err != nil {
		t.Fatalf("LoginOpenAccount() error = %v", err)
	}
}

package twitterscraper

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const bearerToken string = "AAAAAAAAAAAAAAAAAAAAAPYXBAAAAAAACLXUNDekMxqa8h%2F40K4moUkGsoc%3DTYfbDKbT3jJPCEVnMYqilB28NHfOPqkca3qaAxGfsyKCs0wRbw"

// RequestAPI get JSON from frontend API and decodes it
func (s *Scraper) RequestAPI(ctx context.Context, req *http.Request, target interface{}) error {
	s.wg.Wait()
	if s.delay > 0 {
		defer func() {
			s.wg.Add(1)
			go func() {
				time.Sleep(time.Second * time.Duration(s.delay))
				s.wg.Done()
			}()
		}()
	}

	if !s.isLogged {
		if !s.IsGuestToken() || s.guestCreatedAt.Before(time.Now().Add(-time.Hour*3)) {
			err := s.GetGuestToken(ctx)
			if err != nil {
				return err
			}
		}
		req.Header.Set("X-Guest-Token", s.guestToken)
	}

	if s.oAuthToken != "" && s.oAuthSecret != "" {
		req.Header.Set("Authorization", s.sign(req.Method, req.URL))
	} else {
		req.Header.Set("Authorization", "Bearer "+s.bearerToken)
	}

	for _, cookie := range s.client.Jar.Cookies(req.URL) {
		if cookie.Name == "ct0" {
			req.Header.Set("X-CSRF-Token", cookie.Value)
			break
		}
	}

	if s.userAgent != nil {
		req.Header.Set("User-Agent", *s.userAgent)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			return ErrRateLimitExceeded{description: content}
		default:
			return ErrOther{description: content, status: resp.Status, StatusCode: resp.StatusCode}
		}
	}

	if resp.Header.Get("X-Rate-Limit-Remaining") == "0" {
		s.guestToken = ""
	}

	if target == nil {
		return nil
	}
	return jsoniter.NewDecoder(resp.Body).Decode(target)
}

// GetGuestToken from Twitter API
func (s *Scraper) GetGuestToken(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.twitter.com/1.1/guest/activate.json", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.bearerToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			return ErrRateLimitExceeded{description: content}
		default:
			return ErrOther{description: content, status: resp.Status, StatusCode: resp.StatusCode}
		}
	}

	var jsn map[string]interface{}
	if err = jsoniter.NewDecoder(resp.Body).Decode(&jsn); err != nil {
		return errors.Join(err, ErrDecodeGuestTokenResponse)
	}
	var ok bool
	if s.guestToken, ok = jsn["guest_token"].(string); !ok {
		return ErrGuestTokenNotFound
	}
	s.guestCreatedAt = time.Now()

	return nil
}

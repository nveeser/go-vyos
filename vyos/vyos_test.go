package vyos

import (
	"net/http"
	"testing"
)

// TestNewClient tests the NewClient function and its options using t.Run.
func TestNewClient(t *testing.T) {
	baseURL := "https://vyos.example.com"

	t.Run("DefaultClient", func(t *testing.T) {
		c, err := NewClient(baseURL)
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}
		if c.baseURL.String() != baseURL {
			t.Errorf("got baseURL %s, want %s", c.baseURL.String(), baseURL)
		}
		if c.userAgent != defaultUserAgent {
			t.Errorf("got userAgent %s, want %s", c.userAgent, defaultUserAgent)
		}
		if c.httpClient == nil {
			t.Error("httpClient is nil")
		}
	})

	t.Run("WithToken", func(t *testing.T) {
		token := "secret-token"
		c, err := NewClient(baseURL, Token(token))
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}
		if c.token != token {
			t.Errorf("got token %s, want %s", c.token, token)
		}
	})

	t.Run("WithInsecure", func(t *testing.T) {
		c, err := NewClient(baseURL, Insecure())
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}
		tr, ok := c.httpClient.Transport.(*http.Transport)
		if !ok {
			t.Fatal("httpClient transport is not *http.Transport")
		}
		if tr.TLSClientConfig == nil || !tr.TLSClientConfig.InsecureSkipVerify {
			t.Error("InsecureSkipVerify not set")
		}
	})
}

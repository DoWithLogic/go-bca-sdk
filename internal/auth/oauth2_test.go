package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestOAuth2Authenticator_Authenticate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `{
			"access_token": "test-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`)
	}))
	defer server.Close()

	authenticator := NewOAuth2Authenticator("client-id", "client-secret", server.Client(), server.URL)
	req, err := http.NewRequest(http.MethodGet, "https://api.klikbca.com/api/oauth/token", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	err = authenticator.Authenticate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Bearer test-token"
	if got := req.Header.Get("Authorization"); got != expected {
		t.Errorf("expected Authorization %q, got %q", expected, got)
	}
}

func TestOAuth2Authenticator_GetAccessToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method %s, got %s", http.MethodPost, r.Method)
		}

		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("expected Content-Type %q, got %q", "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept %q, got %q", "application/json", r.Header.Get("Accept"))
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			t.Fatal("expected Basic Auth")
		}

		if username != "client-id" {
			t.Errorf("expected client ID %q, got %q", "client-id", username)
		}

		if password != "client-secret" {
			t.Errorf("expected client secret %q, got %q", "client-secret", password)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}

		if got := values.Get("grant_type"); got != "client_credentials" {
			t.Errorf("expected grant_type %q, got %q", "client_credentials", got)
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `{
			"access_token": "test-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`)
	}))
	defer server.Close()

	authenticator := NewOAuth2Authenticator("client-id", "client-secret", server.Client(), server.URL)
	token, err := authenticator.getAccessToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token.AccessToken != "test-token" {
		t.Errorf("expected access token %q, got %q", "test-token", token.AccessToken)
	}

	if token.TokenType != "Bearer" {
		t.Errorf("expected token type %q, got %q", "Bearer", token.TokenType)
	}

	if token.ExpiresIn != 3600 {
		t.Errorf("expected expires_in %d, got %d", 3600, token.ExpiresIn)
	}
}

func TestOAuth2Authenticator_GetAccessToken_NonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid client", http.StatusUnauthorized)
	}))
	defer server.Close()

	authenticator := NewOAuth2Authenticator("client-id", "client-secret", server.Client(), server.URL)
	if _, err := authenticator.getAccessToken(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOAuth2Authenticator_GetAccessToken_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `invalid-json`)
	}))
	defer server.Close()

	authenticator := NewOAuth2Authenticator("client-id", "client-secret", server.Client(), server.URL)
	if _, err := authenticator.getAccessToken(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOAuth2Authenticator_GetAccessToken_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid client", http.StatusUnauthorized)
	}))
	defer server.Close()

	authenticator := NewOAuth2Authenticator("client-id", "client-secret", server.Client(), server.URL)
	_, err := authenticator.getAccessToken(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := "token request failed with status 401 Unauthorized"

	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestOAuth2Authenticator_Authenticate_ReusesValidToken(t *testing.T) {
	var tokenRequests int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenRequests++

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"access_token": "test-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`)
	}))
	defer server.Close()

	authenticator := NewOAuth2Authenticator("client-id", "client-secret", server.Client(), server.URL)

	for range 3 {
		req, err := http.NewRequest(
			http.MethodGet,
			"https://api.klikbca.com/api/oauth/token",
			nil,
		)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		if err := authenticator.Authenticate(context.Background(), req); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("expected Authorization %q, got %q", "Bearer test-token", got)
		}
	}

	if tokenRequests != 1 {
		t.Errorf("expected 1 token request, got %d", tokenRequests)
	}
}

func TestOAuth2Authenticator_Authenticate_RefreshesExpiredToken(t *testing.T) {
	var tokenRequests int
	currentTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenRequests++

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprintf(w, `{
			"access_token": "token-%d",
			"token_type": "Bearer",
			"expires_in": 1
		}`, tokenRequests)
	}))
	defer server.Close()

	authenticator := NewOAuth2Authenticator("client-id", "client-secret", server.Client(), server.URL)
	authenticator.now = func() time.Time {
		return currentTime
	}

	req1, err := http.NewRequest(http.MethodGet, "https://api.klikbca.com/api/oauth/token", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if err := authenticator.Authenticate(context.Background(), req1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req1.Header.Get("Authorization"); got != "Bearer token-1" {
		t.Errorf("expected Authorization %q, got %q", "Bearer token-1", got)
	}

	currentTime = currentTime.Add(1100 * time.Millisecond)

	req2, err := http.NewRequest(http.MethodGet, "https://api.klikbca.com/api/oauth/token", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if err := authenticator.Authenticate(context.Background(), req2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req2.Header.Get("Authorization"); got != "Bearer token-2" {
		t.Errorf("expected Authorization %q, got %q", "Bearer token-2", got)
	}

	if tokenRequests != 2 {
		t.Errorf("expected 2 token requests, got %d", tokenRequests)
	}
}

func TestOAuth2Authenticator_Authenticate_Concurrent(t *testing.T) {
	var tokenRequests int
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		tokenRequests++
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `{
			"access_token": "test-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`)
	}))
	defer server.Close()

	authenticator := NewOAuth2Authenticator("client-id", "client-secret", server.Client(), server.URL)

	const requests = 100

	var wg sync.WaitGroup
	wg.Add(requests)

	for range requests {
		go func() {
			defer wg.Done()

			req, err := http.NewRequest(
				http.MethodGet,
				"https://api.klikbca.com/api/oauth/token",
				nil,
			)
			if err != nil {
				t.Errorf("failed to create request: %v", err)
				return
			}

			if err := authenticator.Authenticate(context.Background(), req); err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
				t.Errorf("expected Authorization %q, got %q", "Bearer test-token", got)
			}
		}()
	}

	wg.Wait()

	if tokenRequests != 1 {
		t.Errorf("expected 1 token request, got %d", tokenRequests)
	}
}

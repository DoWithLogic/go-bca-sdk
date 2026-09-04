package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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

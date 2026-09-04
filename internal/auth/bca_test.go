package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBCAAuthenticator_Authenticate(t *testing.T) {
	var tokenRequests int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenRequests++

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"access_token": "test-access-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`)
	}))
	defer server.Close()

	oauth := NewOAuth2Authenticator("client-id", "client-secret", server.Client(), server.URL)

	authenticator := NewBCAAuthenticator(oauth, "test-api-secret")
	currentTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	authenticator.now = func() time.Time { return currentTime }

	req, err := http.NewRequest(http.MethodPost, "https://api.klikbca.com/banking/corporates/transfers", strings.NewReader(`{"amount":"10000"}`))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	err = authenticator.Authenticate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req.Header.Get("Authorization"); got != "Bearer test-access-token" {
		t.Errorf("expected Authorization %q, got %q", "Bearer test-access-token", got)
	}

	if got := req.Header.Get("X-BCA-Signature"); got == "" {
		t.Error("expected X-BCA-Signature header")
	}

	if tokenRequests != 1 {
		t.Errorf("expected 1 token request, got %d", tokenRequests)
	}

	expectedTimestamp := currentTime.Format(time.RFC3339Nano)
	if got := req.Header.Get("X-BCA-Timestamp"); got != expectedTimestamp {
		t.Errorf("expected X-BCA-Timestamp %q, got %q", expectedTimestamp, got)
	}
}

func TestBCAAuthenticator_Authenticate_RestoresRequestBody(t *testing.T) {
	var tokenRequests int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenRequests++

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"access_token": "test-access-token",
			"token_type": "Bearer",
			"expires_in": 3600
		}`)
	}))
	defer server.Close()

	oauth := NewOAuth2Authenticator("client-id", "client-secret", server.Client(), server.URL)
	authenticator := NewBCAAuthenticator(oauth, "test-api-secret")

	currentTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	authenticator.now = func() time.Time {
		return currentTime
	}

	requestBody := `{"amount":"10000","currency":"IDR"}`
	req, err := http.NewRequest(http.MethodPost, "https://api.klikbca.com/banking/corporates/transfers", strings.NewReader(requestBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if err := authenticator.Authenticate(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("failed to read request body: %v", err)
	}

	if got := string(body); got != requestBody {
		t.Errorf("expected request body %q, got %q", requestBody, got)
	}

	if tokenRequests != 1 {
		t.Errorf("expected 1 token request, got %d", tokenRequests)
	}
}

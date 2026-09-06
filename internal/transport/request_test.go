package transport

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DoWithLogic/go-bca-sdk/errors"
)

type mockAuthenticator struct {
	called bool
}

func (m *mockAuthenticator) Authenticate(ctx context.Context, req *http.Request) error {
	m.called = true
	req.Header.Set("Authorization", "Bearer test-token")
	return nil
}

func TestClient_Do_AuthenticatesRequest(t *testing.T) {
	authenticator := &mockAuthenticator{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("expected Authorization %q, got %q", "Bearer test-token", got)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, authenticator, RetryConfig{MaxRetries: 2, Backoff: 100 * time.Microsecond})

	var result struct {
		Status string `json:"status"`
	}

	err := client.Do(context.Background(), Request{Method: http.MethodGet, Path: "/test"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != "ok" {
		t.Errorf("expected status %q, got %q", "ok", result.Status)
	}

	if !authenticator.called {
		t.Fatal("expected authenticator to be called")
	}
}

func TestClient_Do_NonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		_, _ = w.Write([]byte(`{
			"responseCode": "400US01",
			"responseMessage": "Invalid Field Format cardNumber"
		}`))
	}))
	defer server.Close()

	client := NewClient(
		http.DefaultClient,
		server.URL,
		nil,
		RetryConfig{MaxRetries: 2, Backoff: 100 * time.Microsecond},
	)

	err := client.Do(context.Background(), Request{Method: http.MethodGet, Path: "/balance"}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*errors.APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.HTTPStatusCode != http.StatusBadRequest {
		t.Errorf("expected HTTP status %d, got %d", http.StatusBadRequest, apiErr.HTTPStatusCode)
	}

	if apiErr.ResponseCode != "400US01" {
		t.Errorf("expected response code %q, got %q", "400US01", apiErr.ResponseCode)
	}

	if apiErr.ResponseMessage != "Invalid Field Format cardNumber" {
		t.Errorf("expected response message %q, got %q", "Invalid Field Format cardNumber", apiErr.ResponseMessage)
	}
}

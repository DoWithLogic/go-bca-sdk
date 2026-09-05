package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDefaultRetryPolicy_ShouldRetry(t *testing.T) {
	policy := DefaultRetryPolicy{}

	tests := []struct {
		name       string
		method     string
		statusCode int
		want       bool
	}{
		{"GET server error", http.MethodGet, http.StatusInternalServerError, true},
		{"GET timeout", http.MethodGet, http.StatusRequestTimeout, true},
		{"GET rate limit", http.MethodGet, http.StatusTooManyRequests, true},
		{"POST server error", http.MethodPost, http.StatusInternalServerError, false},
		{"POST bad request", http.MethodPost, http.StatusBadRequest, false},
		{"GET bad request", http.MethodGet, http.StatusBadRequest, false},
		{"DELETE server error", http.MethodDelete, http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := policy.ShouldRetry(tt.method, tt.statusCode)
			if got != tt.want {
				t.Errorf("ShouldRetry(%s, %d) = %v, want %v", tt.method, tt.statusCode, got, tt.want)
			}
		})
	}
}

func TestClient_Do_RetriesGETOnServerError(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"balance":"1000000"}`))
	}))
	defer server.Close()

	client := NewClient(http.DefaultClient, server.URL, nil, RetryConfig{MaxRetries: 2, Backoff: 100 * time.Microsecond})

	var response struct {
		Balance string `json:"balance"`
	}

	err := client.Do(context.Background(), http.MethodGet, "/balance", nil, &response)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.Balance != "1000000" {
		t.Errorf("expected balance 1000000, got %s", response.Balance)
	}
	if requests != 2 {
		t.Errorf("expected 2 requests, got %d", requests)
	}
}

func TestClient_Do_DoesNotRetryPOSTOnServerError(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"server error"}`))
	}))
	defer server.Close()

	client := NewClient(http.DefaultClient, server.URL, nil, RetryConfig{MaxRetries: 2, Backoff: 100 * time.Microsecond})

	err := client.Do(
		context.Background(),
		http.MethodPost,
		"/payment",
		struct {
			Amount int `json:"amount"`
		}{Amount: 100000},
		nil,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if requests != 1 {
		t.Errorf("expected 1 request, got %d", requests)
	}
}

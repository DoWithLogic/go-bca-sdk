package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Do(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}

			w.Header().Set("Content-Type", "application/json")

			_, _ = w.Write([]byte(`{"balance": "1000000"}`))
		}),
	)

	defer server.Close()

	client := NewClient(http.DefaultClient, server.URL, &mockAuthenticator{})

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
}

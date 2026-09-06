package transport

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Do_APIError(t *testing.T) {
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
		server.Client(),
		server.URL,
		nil,
		RetryConfig{},
	)

	err := client.Do(
		context.Background(),
		Request{
			Method: http.MethodPost,
			Path:   "/openapi/bdc/v1.0/update-status",
			Body: map[string]string{
				"cardNumber": "123",
			},
		},
		nil,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.HTTPStatusCode != http.StatusBadRequest {
		t.Errorf(
			"expected HTTP status %d, got %d",
			http.StatusBadRequest,
			apiErr.HTTPStatusCode,
		)
	}

	if apiErr.ResponseCode != "400US01" {
		t.Errorf(
			"expected response code %q, got %q",
			"400US01",
			apiErr.ResponseCode,
		)
	}

	if apiErr.ResponseMessage != "Invalid Field Format cardNumber" {
		t.Errorf(
			"expected response message %q, got %q",
			"Invalid Field Format cardNumber",
			apiErr.ResponseMessage,
		)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		HTTPStatusCode:  http.StatusBadRequest,
		ResponseCode:    "400US01",
		ResponseMessage: "Invalid Field Format cardNumber",
	}

	got := err.Error()

	want := "BCA API error: 400US01 - Invalid Field Format cardNumber"

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

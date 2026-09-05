package transport

import (
	"net/http"
	"testing"
)

func TestAPIError_Unmarshal(t *testing.T) {
	apiErr := &APIError{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte(`{"error":"invalid request"}`),
	}

	var response struct {
		Error string `json:"error"`
	}

	if err := apiErr.Unmarshal(&response); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Error != "invalid request" {
		t.Errorf("expected error %q, got %q", "invalid request", response.Error)
	}
}

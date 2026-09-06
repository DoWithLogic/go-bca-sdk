package errors

import (
	"net/http"
	"testing"
)

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

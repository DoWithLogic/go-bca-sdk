package errors

import (
	"fmt"
)

// APIError represents an error returned by the BCA API.
type APIError struct {
	HTTPStatusCode  int
	ResponseCode    string
	ResponseMessage string
}

// Error returns a human-readable representation of the API error.
func (e *APIError) Error() string {
	return fmt.Sprintf("BCA API error: %s - %s", e.ResponseCode, e.ResponseMessage)
}

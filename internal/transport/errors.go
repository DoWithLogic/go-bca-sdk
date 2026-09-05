package transport

import (
	"encoding/json"
	"fmt"
)

// APIError represents an error returned by the BCA API.
type APIError struct {
	StatusCode int
	Status     string
	Body       []byte
}

// Error returns a human-readable representation of the API error.
func (e *APIError) Error() string {
	return fmt.Sprintf("bca api request failed with status %s", e.Status)
}

// UnmarshalJSON decodes the API error response body into v.
func (e *APIError) Unmarshal(v any) error { return json.Unmarshal(e.Body, v) }

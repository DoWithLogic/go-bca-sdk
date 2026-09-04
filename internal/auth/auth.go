package auth

import (
	"context"
	"net/http"
)

// Authenticator authenticates HTTP requests before they are sent to the BCA API.
type Authenticator interface {
	Authenticate(ctx context.Context, req *http.Request) error
}

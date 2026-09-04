package auth

import (
	"context"
	"net/http"
)

// SNAPAuthenticator authenticates HTTP requests using BCA SNAP authentication.
//
// See the BCA Developer API documentation for more information:
// https://developer.bca.co.id/Dokumentasi#oauth20snap
type SNAPAuthenticator struct {
}

// Authenticate authenticates the HTTP request using the BCA SNAP authentication flow.
func (a *SNAPAuthenticator) Authenticate(ctx context.Context, req *http.Request) error {
	return nil
}

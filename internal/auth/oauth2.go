package auth

import (
	"context"
	"net/http"
)

// OAuth2Authenticator authenticates HTTP requests using BCA OAuth 2.0 authentication.
//
// See the BCA Developer API documentation for more information:
// https://developer.bca.co.id/Dokumentasi#oauth20
type OAuth2Authenticator struct {
}

// Authenticate authenticates the HTTP request using the BCA OAuth 2.0 authentication flow.
func (a *OAuth2Authenticator) Authenticate(ctx context.Context, req *http.Request) error {
	return nil
}

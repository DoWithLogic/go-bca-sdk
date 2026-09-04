package auth

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/DoWithLogic/go-bca-sdk/internal/signature"
)

// BCAAuthenticator authenticates HTTP requests using OAuth 2.0 and BCA request signing.
type BCAAuthenticator struct {
	oauth     *OAuth2Authenticator
	apiSecret string
	now       func() time.Time
}

// NewBCAAuthenticator creates a new BCA authenticator using OAuth 2.0 and request signing.
func NewBCAAuthenticator(oauth *OAuth2Authenticator, apiSecret string) *BCAAuthenticator {
	return &BCAAuthenticator{
		oauth:     oauth,
		apiSecret: apiSecret,
		now:       time.Now,
	}
}

// Authenticate authenticates the HTTP request using OAuth 2.0 and BCA request signing.
func (a *BCAAuthenticator) Authenticate(ctx context.Context, req *http.Request) error {
	token, err := a.oauth.getToken(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token.tokenType+" "+token.accessToken)

	var body string
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		body = string(bodyBytes)
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	timestamp := a.now().Format(time.RFC3339Nano)
	req.Header.Set("X-BCA-Timestamp", timestamp)

	sig := signature.Sign(req.Method, req.URL.RequestURI(), token.accessToken, body, timestamp, a.apiSecret)
	req.Header.Set("X-BCA-Signature", sig)

	return nil
}

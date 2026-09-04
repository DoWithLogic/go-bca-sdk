package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// OAuth2Authenticator authenticates HTTP requests using BCA OAuth 2.0 authentication.
//
// See the BCA Developer API documentation for more information:
// https://developer.bca.co.id/Dokumentasi#oauth20
type OAuth2Authenticator struct {
	clientID     string
	clientSecret string
	tokenURL     string
	httpClient   *http.Client

	now   func() time.Time
	mu    sync.Mutex
	token *token
}

func NewOAuth2Authenticator(clientID, clientSecret string, httpClient *http.Client, tokenURL string) *OAuth2Authenticator {
	return &OAuth2Authenticator{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   httpClient,
		tokenURL:     tokenURL,
		now:          time.Now,
	}
}

// Authenticate authenticates the HTTP request using the BCA OAuth 2.0 authentication flow.
func (a *OAuth2Authenticator) Authenticate(ctx context.Context, req *http.Request) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.token != nil && a.now().Before(a.token.expiresAt) {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", a.token.tokenType, a.token.accessToken))

		return nil
	}

	tokenResponse, err := a.getAccessToken(ctx)
	if err != nil {
		return err
	}

	a.token = &token{
		accessToken: tokenResponse.AccessToken,
		tokenType:   tokenResponse.TokenType,
		expiresAt:   a.now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenResponse.TokenType, tokenResponse.AccessToken))

	return nil
}

func (a *OAuth2Authenticator) getAccessToken(ctx context.Context) (*tokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.tokenURL, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(a.clientID, a.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("token request failed with status %s", resp.Status)
	}

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type token struct {
	accessToken string
	tokenType   string
	expiresAt   time.Time
}

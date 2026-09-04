package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
}

func NewOAuth2Authenticator(clientID, clientSecret string, httpClient *http.Client, tokenURL string) *OAuth2Authenticator {
	return &OAuth2Authenticator{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   httpClient,
		tokenURL:     tokenURL,
	}
}

// Authenticate authenticates the HTTP request using the BCA OAuth 2.0 authentication flow.
func (a *OAuth2Authenticator) Authenticate(ctx context.Context, req *http.Request) error {
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

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

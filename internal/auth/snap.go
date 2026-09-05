package auth

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// SNAPAuthenticator authenticates HTTP requests using BCA SNAP authentication.
//
// See the BCA Developer API documentation for more information:
// https://developer.bca.co.id/Dokumentasi#oauth20snap
type SNAPAuthenticator struct {
	clientID   string
	privateKey *rsa.PrivateKey
	httpClient *http.Client
	tokenURL   string
	now        func() time.Time
}

type snapTokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   string `json:"expiresIn"`
}

func NewSNAPAuthenticator(clientID string, privateKey *rsa.PrivateKey, httpClient *http.Client, tokenURL string) *SNAPAuthenticator {
	return &SNAPAuthenticator{
		clientID:   clientID,
		privateKey: privateKey,
		httpClient: httpClient,
		tokenURL:   tokenURL,
		now:        time.Now,
	}
}

// Authenticate authenticates the HTTP request using the BCA SNAP authentication flow.
func (a *SNAPAuthenticator) Authenticate(ctx context.Context, req *http.Request) error {
	return nil
}

func generateSNAPSignature(clientID string, timestamp string, privateKey *rsa.PrivateKey) (string, error) {
	stringToSign := clientID + "|" + timestamp
	hash := sha256.Sum256([]byte(stringToSign))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func (a *SNAPAuthenticator) getAccessToken(ctx context.Context) (*snapTokenResponse, error) {
	timestamp := a.now().Format(time.RFC3339)

	signature, err := generateSNAPSignature(a.clientID, timestamp, a.privateKey)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.tokenURL, strings.NewReader(`{"grantType":"client_credentials"}`))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TIMESTAMP", timestamp)
	req.Header.Set("X-CLIENT-KEY", a.clientID)
	req.Header.Set("X-SIGNATURE", signature)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("SNAP token request failed with status %s", resp.Status)
	}

	var token snapTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

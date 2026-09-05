package auth

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DoWithLogic/go-bca-sdk/internal/signature"
)

// SNAPAuthenticator authenticates HTTP requests using BCA SNAP authentication.
//
// See the BCA Developer API documentation for more information:
// https://developer.bca.co.id/Dokumentasi#oauth20snap
type SNAPAuthenticator struct {
	clientID     string
	clientSecret string
	privateKey   *rsa.PrivateKey
	httpClient   *http.Client
	tokenURL     string
	now          func() time.Time

	mu    sync.Mutex
	token *snapToken
}

type snapTokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   string `json:"expiresIn"`
}

type snapToken struct {
	AccessToken string
	TokenType   string
	ExpiresAt   time.Time
}

// NewSNAPAuthenticator creates a new SNAP authenticator.
func NewSNAPAuthenticator(clientID, clientSecret string, privateKey *rsa.PrivateKey, httpClient *http.Client, tokenURL string) *SNAPAuthenticator {
	return &SNAPAuthenticator{
		clientID:     clientID,
		clientSecret: clientSecret,
		privateKey:   privateKey,
		httpClient:   httpClient,
		tokenURL:     tokenURL,
		now:          time.Now,
	}
}

// Authenticate authenticates the HTTP request using the BCA SNAP authentication flow.
func (a *SNAPAuthenticator) Authenticate(ctx context.Context, req *http.Request) error {
	token, err := a.getToken(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token.TokenType+" "+token.AccessToken)

	timestamp := a.now().Format(time.RFC3339)
	req.Header.Set("X-TIMESTAMP", timestamp)

	var body string
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}

		body = string(bodyBytes)

		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	sig, err := signature.SignSNAP(req.Method, req.URL.RequestURI(), token.AccessToken, body, timestamp, a.clientSecret)
	if err != nil {
		return err
	}

	req.Header.Set("X-SIGNATURE", sig)

	return nil
}

// generateSNAPSignature generates an RSA SHA-256 signature for the BCA SNAP
// access-token request.
//
// The signature is generated from:
//
//	ClientID|Timestamp
//
// The resulting RSA signature is Base64-encoded.
func generateSNAPSignature(clientID string, timestamp string, privateKey *rsa.PrivateKey) (string, error) {
	stringToSign := clientID + "|" + timestamp
	hash := sha256.Sum256([]byte(stringToSign))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// getAccessToken requests a new SNAP access token from the BCA token endpoint.
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

// getToken returns a cached SNAP access token or requests a new token
// when the cached token is missing or expired.
func (a *SNAPAuthenticator) getToken(ctx context.Context) (*snapToken, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.token != nil && a.now().Before(a.token.ExpiresAt) {
		return a.token, nil
	}

	response, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	expiresIn, err := strconv.Atoi(response.ExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("invalid SNAP token expiry : %w", err)
	}

	a.token = &snapToken{
		AccessToken: response.AccessToken,
		TokenType:   response.TokenType,
		ExpiresAt:   a.now().Add(time.Duration(expiresIn) * time.Second),
	}

	return a.token, nil
}

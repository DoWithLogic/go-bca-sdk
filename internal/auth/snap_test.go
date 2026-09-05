package auth

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestGenerateSNAPSignature(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	clientID := "client-id"
	timestamp := "2026-09-05T12:00:00+07:00"

	signature, err := generateSNAPSignature(
		clientID,
		timestamp,
		privateKey,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if signature == "" {
		t.Fatal("expected signature, got empty string")
	}

	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		t.Fatalf("signature is not valid base64: %v", err)
	}

	hash := sha256.Sum256([]byte(clientID + "|" + timestamp))

	if err := rsa.VerifyPKCS1v15(
		&privateKey.PublicKey,
		crypto.SHA256,
		hash[:],
		decoded,
	); err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}
}

func TestSNAPAuthenticator_GetAccessToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	authenticator := NewSNAPAuthenticator("client-id", "client-secret", privateKey, nil, "")
	authenticator.now = func() time.Time {
		return time.Date(2026, 9, 5, 12, 0, 0, 0, time.FixedZone("WIB", 7*60*60))
	}

	var receivedSignature string
	var receivedClientID string
	var receivedTimestamp string
	var receivedBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/openapi/v1.0/access-token/b2b" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", got)
		}

		receivedClientID = r.Header.Get("X-CLIENT-KEY")
		receivedTimestamp = r.Header.Get("X-TIMESTAMP")
		receivedSignature = r.Header.Get("X-SIGNATURE")

		receivedBody, _ = io.ReadAll(r.Body)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"accessToken": "test-access-token",
			"tokenType": "Bearer",
			"expiresIn": "900"
		}`))
	}))
	defer server.Close()

	authenticator.httpClient = server.Client()
	authenticator.tokenURL = server.URL + "/openapi/v1.0/access-token/b2b"

	token, err := authenticator.getAccessToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token.AccessToken != "test-access-token" {
		t.Errorf("expected access token %q, got %q", "test-access-token", token.AccessToken)
	}

	if token.TokenType != "Bearer" {
		t.Errorf("expected token type %q, got %q", "Bearer", token.TokenType)
	}

	if token.ExpiresIn != "900" {
		t.Errorf("expected expires in %q, got %q", "900", token.ExpiresIn)
	}

	if receivedClientID != "client-id" {
		t.Errorf("expected X-CLIENT-KEY %q, got %q", "client-id", receivedClientID)
	}

	expectedTimestamp := "2026-09-05T12:00:00+07:00"

	if receivedTimestamp != expectedTimestamp {
		t.Errorf("expected X-TIMESTAMP %q, got %q", expectedTimestamp, receivedTimestamp)
	}

	if string(receivedBody) != `{"grantType":"client_credentials"}` {
		t.Errorf("unexpected request body: %s", receivedBody)
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(receivedSignature)
	if err != nil {
		t.Fatalf("invalid base64 signature: %v", err)
	}

	hash := sha256.Sum256(
		[]byte("client-id|" + expectedTimestamp),
	)

	err = rsa.VerifyPKCS1v15(&privateKey.PublicKey, crypto.SHA256, hash[:], decodedSignature)
	if err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}
}

func TestSNAPAuthenticator_GetAccessToken_NonSuccess(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid client", http.StatusUnauthorized)
	}))
	defer server.Close()

	authenticator := NewSNAPAuthenticator("client-id", "client-secret", privateKey, server.Client(), server.URL)

	if _, err = authenticator.getAccessToken(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSNAPAuthenticator_GetAccessToken_InvalidJSON(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`invalid-json`))
	}))
	defer server.Close()

	authenticator := NewSNAPAuthenticator("client-id", "client-secret", privateKey, server.Client(), server.URL)

	if _, err = authenticator.getAccessToken(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSNAPAuthenticator_GetToken_ReusesValidToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"accessToken": "test-access-token",
			"tokenType": "Bearer",
			"expiresIn": "900"
		}`))
	}))
	defer server.Close()

	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.FixedZone("WIB", 7*60*60))
	authenticator := NewSNAPAuthenticator("client-id", "client-secret", privateKey, server.Client(), server.URL)
	authenticator.now = func() time.Time { return now }

	token1, err := authenticator.getToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	token2, err := authenticator.getToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token1.AccessToken != token2.AccessToken {
		t.Errorf("expected same access token, got %q and %q", token1.AccessToken, token2.AccessToken)
	}

	if requestCount != 1 {
		t.Errorf("expected 1 token request, got %d", requestCount)
	}
}

func TestSNAPAuthenticator_GetToken_RefreshesExpiredToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		w.Header().Set("Content-Type", "application/json")

		if requestCount == 1 {
			_, _ = w.Write([]byte(`{
				"accessToken": "token-1",
				"tokenType": "Bearer",
				"expiresIn": "900"
			}`))
			return
		}

		_, _ = w.Write([]byte(`{
			"accessToken": "token-2",
			"tokenType": "Bearer",
			"expiresIn": "900"
		}`))
	}))
	defer server.Close()

	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.FixedZone("WIB", 7*60*60))
	authenticator := NewSNAPAuthenticator("client-id", "client-secret", privateKey, server.Client(), server.URL)
	authenticator.now = func() time.Time { return now }

	token1, err := authenticator.getToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token1.AccessToken != "token-1" {
		t.Errorf("expected token-1, got %q", token1.AccessToken)
	}

	// Move time beyond the 15-minute expiration.
	now = now.Add(16 * time.Minute)

	token2, err := authenticator.getToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token2.AccessToken != "token-2" {
		t.Errorf("expected token-2, got %q", token2.AccessToken)
	}

	if requestCount != 2 {
		t.Errorf("expected 2 token requests, got %d", requestCount)
	}
}

func TestSNAPAuthenticator_GetToken_Concurrent(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"accessToken": "test-access-token",
			"tokenType": "Bearer",
			"expiresIn": "900"
		}`))
	}))
	defer server.Close()

	authenticator := NewSNAPAuthenticator("client-id", "client-secret", privateKey, server.Client(), server.URL)
	authenticator.now = func() time.Time {
		return time.Date(2026, 9, 5, 12, 0, 0, 0, time.FixedZone("WIB", 7*60*60))
	}

	const goroutines = 20

	var wg sync.WaitGroup
	wg.Add(goroutines)

	errors := make(chan error, goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()

			token, err := authenticator.getToken(context.Background())
			if err != nil {
				errors <- err
				return
			}

			if token.AccessToken != "test-access-token" {
				errors <- fmt.Errorf(
					"unexpected access token: %s",
					token.AccessToken,
				)
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}

	if requestCount != 1 {
		t.Errorf("expected 1 token request, got %d", requestCount)
	}
}

func TestSNAPAuthenticator_Authenticate(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"accessToken": "test-access-token",
			"tokenType": "Bearer",
			"expiresIn": "900"
		}`))
	}))
	defer server.Close()

	authenticator := NewSNAPAuthenticator("client-id", "client-secret", privateKey, server.Client(), server.URL)

	req, err := http.NewRequest(http.MethodPost, "https://example.com/openapi/v1.0/balance-inquiry", strings.NewReader(`{"accountNo":"1234567890"}`))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	authenticator.now = func() time.Time {
		return time.Date(2026, 9, 5, 12, 0, 0, 0, time.FixedZone("WIB", 7*60*60))
	}

	err = authenticator.Authenticate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := req.Header.Get("Authorization"); got != "Bearer test-access-token" {
		t.Errorf("unexpected Authorization header: %q", got)
	}

	if got := req.Header.Get("X-TIMESTAMP"); got == "" {
		t.Error("expected X-TIMESTAMP header")
	}

	if got := req.Header.Get("X-SIGNATURE"); got == "" {
		t.Error("expected X-SIGNATURE header")
	}
}

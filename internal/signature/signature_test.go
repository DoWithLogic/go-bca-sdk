package signature

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestSign(t *testing.T) {
	method := "POST"
	relativeURL := "/banking/corporates/transfers"
	accessToken := "lIWOt2p29grUo59bedBUrBY3pnzqQX544LzYPohcGHOuwn8AUEdUKS"
	requestBody := `{"amount":"10000","currency":"IDR"}`
	timestamp := "2016-02-03T10:00:00.000+07:00"
	apiSecret := "22a2d25e-765d-41e1-8d29-da68dcb5698b"

	got := Sign(method, relativeURL, accessToken, requestBody, timestamp, apiSecret)
	if got == "" {
		t.Fatal("expected signature, got empty string")
	}
}

func TestSign_IsDeterministic(t *testing.T) {
	method := "POST"
	relativeURL := "/banking/corporates/transfers"
	accessToken := "test-access-token"
	requestBody := `{"amount":"10000","currency":"IDR"}`
	timestamp := "2026-01-01T12:00:00.000+07:00"
	apiSecret := "test-api-secret"

	got1 := Sign(method, relativeURL, accessToken, requestBody, timestamp, apiSecret)
	got2 := Sign(method, relativeURL, accessToken, requestBody, timestamp, apiSecret)

	if got1 != got2 {
		t.Errorf("expected deterministic signature, got %q and %q", got1, got2)
	}
}

func TestSign_ChangesWithRequestBody(t *testing.T) {
	method := "POST"
	relativeURL := "/banking/corporates/transfers"
	accessToken := "test-access-token"
	timestamp := "2026-01-01T12:00:00.000+07:00"
	apiSecret := "test-api-secret"

	signature1 := Sign(method, relativeURL, accessToken, `{"amount":"10000"}`, timestamp, apiSecret)
	signature2 := Sign(method, relativeURL, accessToken, `{"amount":"20000"}`, timestamp, apiSecret)

	if signature1 == signature2 {
		t.Fatal("expected different signatures for different request bodies")
	}
}

func TestSign_ChangesWithTimestamp(t *testing.T) {
	method := "POST"
	relativeURL := "/banking/corporates/transfers"
	accessToken := "test-access-token"
	requestBody := `{"amount":"10000"}`
	apiSecret := "test-api-secret"

	signature1 := Sign(method, relativeURL, accessToken, requestBody, "2026-01-01T12:00:00.000+07:00", apiSecret)
	signature2 := Sign(method, relativeURL, accessToken, requestBody, "2026-01-01T12:01:00.000+07:00", apiSecret)

	if signature1 == signature2 {
		t.Fatal("expected different signatures for different timestamps")
	}
}

func TestSign_EmptyRequestBody(t *testing.T) {
	method := "GET"
	relativeURL := "/banking/corporates/accounts"
	accessToken := "test-access-token"
	timestamp := "2026-01-01T12:00:00.000+07:00"
	apiSecret := "test-api-secret"

	got := Sign(method, relativeURL, accessToken, "", timestamp, apiSecret)
	if got == "" {
		t.Fatal("expected signature, got empty string")
	}
}

func TestBuildStringToSign(t *testing.T) {
	method := "POST"
	relativeURL := "/banking/corporates/transfers"
	accessToken := "test-access-token"
	requestBody := `{"amount":"10000"}`
	timestamp := "2026-01-01T12:00:00.000+07:00"

	got := buildStringToSign(method, relativeURL, accessToken, requestBody, timestamp)

	bodyHash := sha256.Sum256([]byte(requestBody))
	bodyHashHex := hex.EncodeToString(bodyHash[:])

	expected := method + ":" + relativeURL + ":" + accessToken + ":" + bodyHashHex + ":" + timestamp
	if got != expected {
		t.Errorf("expected StringToSign %q, got %q", expected, got)
	}
}

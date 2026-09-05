package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestSignSNAP(t *testing.T) {
	method := "POST"
	relativeURL := "/openapi/v1.0/balance-inquiry"
	accessToken := "test-access-token"
	requestBody := `{"accountNo":"1234567890"}`
	timestamp := "2026-09-05T12:00:00+07:00"
	clientSecret := "client-secret"

	bodyHash := sha256.Sum256([]byte(requestBody))
	bodyHashHex := hex.EncodeToString(bodyHash[:])

	stringToSign := method + ":" + relativeURL + ":" + accessToken + ":" + bodyHashHex + ":" + timestamp

	mac := hmac.New(sha512.New, []byte(clientSecret))
	mac.Write([]byte(stringToSign))

	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	got, err := SignSNAP(method, relativeURL, accessToken, requestBody, timestamp, clientSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != expected {
		t.Errorf("unexpected signature: expected %q, got %q", expected, got)
	}
}

func TestSignSNAP_MinifiesJSON(t *testing.T) {
	signature1, err := SignSNAP(
		"POST",
		"/openapi/v1.0/balance-inquiry",
		"test-access-token",
		`{"accountNo":"1234567890"}`,
		"2026-09-05T12:00:00+07:00",
		"client-secret",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	signature2, err := SignSNAP(
		"POST",
		"/openapi/v1.0/balance-inquiry",
		"test-access-token",
		`{
			"accountNo": "1234567890"
		}`,
		"2026-09-05T12:00:00+07:00",
		"client-secret",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if signature1 != signature2 {
		t.Errorf("expected signatures to match for equivalent JSON: %q != %q", signature1, signature2)
	}
}

func TestSignSNAP_InvalidJSON(t *testing.T) {
	_, err := SignSNAP(
		"POST",
		"/openapi/v1.0/balance-inquiry",
		"test-access-token",
		`invalid-json`,
		"2026-09-05T12:00:00+07:00",
		"client-secret",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSignSNAP_SortsQueryParameters(t *testing.T) {
	timestamp := "2026-09-05T12:00:00+07:00"

	signature1, err := SignSNAP(
		"GET",
		"/api/test?b=2&a=1",
		"test-access-token",
		"",
		timestamp,
		"client-secret",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	signature2, err := SignSNAP(
		"GET",
		"/api/test?a=1&b=2",
		"test-access-token",
		"",
		timestamp,
		"client-secret",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if signature1 != signature2 {
		t.Errorf(
			"expected signatures to match: %q != %q",
			signature1,
			signature2,
		)
	}
}

func TestCanonicalizeURL(t *testing.T) {
	got, err := canonicalizeURL("/api/test?b=2&a=1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "/api/test?a=1&b=2"

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

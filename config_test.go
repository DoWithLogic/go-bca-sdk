package bca

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/DoWithLogic/go-bca-sdk/internal/auth"
)

func TestConfig_NewAuthenticator_BCA(t *testing.T) {
	cfg := defaultConfig()

	authenticator, err := cfg.newAuthenticator()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := authenticator.(*auth.BCAAuthenticator); !ok {
		t.Fatalf("expected BCAAuthenticator, got %T", authenticator)
	}
}

func TestConfig_NewAuthenticator_SNAP(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	cfg := defaultConfig()
	cfg.AuthMode = AuthModeSNAP
	cfg.SNAPPrivateKey = privateKey

	authenticator, err := cfg.newAuthenticator()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := authenticator.(*auth.SNAPAuthenticator); !ok {
		t.Fatalf("expected SNAPAuthenticator, got %T", authenticator)
	}
}

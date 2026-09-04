package bca

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(
		WithOAuthCredentials("test-client", "test-secret"),
		WithEnvironment(Production),
		WithTimeout(10*time.Second),
		WithAPISecret("api-secret"),
	)

	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.config.ClientID != "test-client" {
		t.Errorf("unexpected client ID")
	}

	if client.config.Environment != Production {
		t.Errorf("unexpected environment")
	}

	if client.config.Timeout != 10*time.Second {
		t.Errorf("unexpected timeout")
	}

	if client.config.APISecret != "api-secret" {
		t.Errorf("unexpected API Secret")
	}
}

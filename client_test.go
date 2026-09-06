package bca

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(
		WithClientID("client-id"),
		WithClientSecret("client-secret"),
		WithEnvironment(Production),
		WithAPISecret("api-secret"),
	)

	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client.config.ClientID != "client-id" {
		t.Errorf("unexpected client ID")
	}

	if client.config.ClientSecret != "client-secret" {
		t.Errorf("unexpected client Secret")
	}

	if client.config.Environment != Production {
		t.Errorf("unexpected environment")
	}

	if client.config.APISecret != "api-secret" {
		t.Errorf("unexpected API Secret")
	}
}

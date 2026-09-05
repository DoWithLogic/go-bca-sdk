package bca

import (
	"testing"
)

func TestWithMaxRetries_Invalid(t *testing.T) {
	cfg := defaultConfig()

	err := WithMaxRetries(-1)(&cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWithRetryBackoff_Invalid(t *testing.T) {
	cfg := defaultConfig()

	err := WithRetryBackoff(-1)(&cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

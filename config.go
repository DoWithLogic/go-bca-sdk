package bca

import (
	"net/http"
	"time"
)

// Environment represents the BCA API environment.
type Environment string

const (
	// Sandbox is the BCA sandbox environment used for testing.
	Sandbox Environment = "sandbox"

	// Production is the BCA production environment used for live requests.
	Production Environment = "production"
)

func (e Environment) is(val Environment) bool { return e == val }
func (e Environment) baseURL() string {
	if !e.is(Sandbox) {
		return "https://api.klikbca.com"
	}

	return "https://sandbox.bca.co.id"
}

// Config contains the configuration used by the BCA client.
type Config struct {
	// Environment specifies the BCA API environment.
	// Defaults to Sandbox.
	Environment Environment

	// BaseURL specifies the base URL for BCA API requests.
	// Defaults to the environment's base URL.
	BaseURL string

	// ClientID is the OAuth 2.0 client ID.
	ClientID string

	// ClientSecret is the OAuth 2.0 client secret.
	ClientSecret string

	// APISecret is the secret used to generate BCA API request signatures.
	APISecret string

	// HTTPClient specifies the HTTP client used for API requests.
	// Defaults to an HTTP client with a 30-second timeout.
	HTTPClient *http.Client

	// MaxRetries specifies the maximum number of retries for retryable requests.
	// Defaults to 2.
	MaxRetries int

	// RetryBackoff specifies the initial delay between retries.
	// Defaults to 100 milliseconds.
	RetryBackoff time.Duration
}

// defaultConfig returns the default configuration for the BCA client.
func defaultConfig() Config {
	return Config{
		Environment:  Sandbox,
		BaseURL:      Sandbox.baseURL(),
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		MaxRetries:   2,
		RetryBackoff: 100 * time.Millisecond,
	}
}

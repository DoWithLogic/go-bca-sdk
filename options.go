package bca

import (
	"fmt"
	"net/http"
	"time"
)

// Option configures a Client.
type Option func(*Config) error

// WithEnvironment sets the BCA API environment to use.
func WithEnvironment(env Environment) Option {
	return func(cfg *Config) error {
		cfg.Environment = env
		cfg.BaseURL = env.baseURL()
		return nil
	}
}

// WithOAuthCredentials sets the OAuth client ID and client secret
// used to obtain an access token from BCA.
func WithOAuthCredentials(clientID, clientSecret string) Option {
	return func(cfg *Config) error {
		cfg.ClientID = clientID
		cfg.ClientSecret = clientSecret
		return nil
	}
}

// WithHTTPClient sets the HTTP client used to make requests to BCA.
func WithHTTPClient(client *http.Client) Option {
	return func(cfg *Config) error {
		cfg.HTTPClient = client
		return nil
	}
}

// WithAPISecret sets the API secret used to sign BCA API requests.
func WithAPISecret(apiSecret string) Option {
	return func(cfg *Config) error {
		cfg.APISecret = apiSecret
		return nil
	}
}

// WithMaxRetries sets the maximum number of retries for retryable requests.
func WithMaxRetries(maxRetries int) Option {
	return func(cfg *Config) error {
		if maxRetries < 0 {
			return fmt.Errorf("max retries cannot be negative")
		}
		cfg.MaxRetries = maxRetries
		return nil
	}
}

// WithRetryBackoff sets the initial backoff duration between retries.
func WithRetryBackoff(backoff time.Duration) Option {
	return func(cfg *Config) error {
		if backoff < 0 {
			return fmt.Errorf("retry backoff cannot be negative")
		}
		cfg.RetryBackoff = backoff
		return nil
	}
}

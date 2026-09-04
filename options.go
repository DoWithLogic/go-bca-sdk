package bca

import (
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

// WithTimeout sets the timeout for HTTP requests made by the client.
func WithTimeout(timeout time.Duration) Option {
	return func(cfg *Config) error {
		cfg.Timeout = timeout
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

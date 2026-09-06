package bca

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"
)

// Option configures a Client.
type Option func(*Config) error

// WithEnvironment sets the BCA API environment to use.
func WithEnvironment(env Environment) Option {
	return func(cfg *Config) error {
		baseURL, err := env.baseURL()
		if err != nil {
			return err
		}

		cfg.Environment = env
		cfg.BaseURL = baseURL
		return nil
	}
}

// WithClientID sets the BCA client ID.
func WithClientID(clientID string) Option {
	return func(cfg *Config) error {
		cfg.ClientID = clientID
		return nil
	}
}

// WithClientSecret sets the BCA client secret.
func WithClientSecret(clientSecret string) Option {
	return func(cfg *Config) error {
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

// WithSNAPAuth configures the client to use BCA SNAP authentication
// with the provided RSA private key.
func WithSNAPAuth(privateKey *rsa.PrivateKey) Option {
	return func(cfg *Config) error {
		cfg.AuthMode = AuthModeSNAP
		cfg.SNAPPrivateKey = privateKey
		return nil
	}
}

// WithChannelID sets the channel identifier.
func WithChannelID(channelID string) Option {
	return func(cfg *Config) error {
		cfg.ChannelID = channelID
		return nil
	}
}

// WithPartnerID sets the Klik BCA Bisnis's Corporate ID
func WithPartnerID(partnerID string) Option {
	return func(cfg *Config) error {
		cfg.PartnerID = partnerID
		return nil
	}
}

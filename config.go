package bca

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/DoWithLogic/go-bca-sdk/internal/auth"
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

// AuthMode specifies the authentication mechanism used by the BCA client.
type AuthMode string

const (
	// AuthModeBCA uses the traditional BCA API authentication.
	AuthModeBCA AuthMode = "bca"

	// AuthModeSNAP uses the BCA SNAP authentication.
	AuthModeSNAP AuthMode = "snap"
)

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

	// AuthMode specifies the authentication mechanism used by the client.
	AuthMode AuthMode

	// SNAPPrivateKey is the RSA private key used to sign SNAP access-token requests.
	SNAPPrivateKey *rsa.PrivateKey
}

// defaultConfig returns the default configuration for the BCA client.
func defaultConfig() Config {
	return Config{
		Environment:  Sandbox,
		BaseURL:      Sandbox.baseURL(),
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		MaxRetries:   2,
		RetryBackoff: 100 * time.Millisecond,
		AuthMode:     AuthModeBCA,
	}
}

// newAuthenticator creates an authenticator based on the configured
// authentication mode.
func (c Config) newAuthenticator() (auth.Authenticator, error) {
	switch c.AuthMode {
	case AuthModeSNAP:
		return auth.NewSNAPAuthenticator(
			c.ClientID,
			c.ClientSecret,
			c.SNAPPrivateKey,
			c.HTTPClient,
			c.BaseURL+"/openapi/v1.0/access-token/b2b",
		), nil
	case AuthModeBCA:
		return auth.NewBCAAuthenticator(
			auth.NewOAuth2Authenticator(c.ClientID, c.ClientSecret, c.HTTPClient, c.BaseURL+"/api/oauth/token"),
			c.APISecret,
		), nil
	default:
		return nil, fmt.Errorf("unsupported auth mode: %q", c.AuthMode)
	}
}

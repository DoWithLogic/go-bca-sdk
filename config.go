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
	Environment Environment
	BaseURL     string

	ClientID     string
	ClientSecret string
	APISecret    string

	HTTPClient *http.Client
	Timeout    time.Duration
}

// defaultConfig returns the default configuration for the BCA client.
func defaultConfig() Config {
	return Config{
		Environment: Sandbox,
		BaseURL:     Sandbox.baseURL(),
		HTTPClient:  &http.Client{Timeout: 30 * time.Second},
		Timeout:     30 * time.Second,
	}
}

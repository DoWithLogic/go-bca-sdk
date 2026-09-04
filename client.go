package bca

import (
	"github.com/DoWithLogic/go-bca-sdk/internal/auth"
	"github.com/DoWithLogic/go-bca-sdk/internal/transport"
)

// Client is the main client for interacting with the BCA API.
type Client struct {
	config    Config
	transport *transport.Client
}

// NewClient creates a new BCA API client using the provided options.
func NewClient(opts ...Option) (*Client, error) {
	cfg := defaultConfig()

	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	oauth := auth.NewOAuth2Authenticator(cfg.ClientID, cfg.ClientSecret, cfg.HTTPClient, cfg.BaseURL+"/api/oauth/token")
	authenticator := auth.NewBCAAuthenticator(oauth, cfg.APISecret)

	client := &Client{
		config:    cfg,
		transport: transport.NewClient(cfg.HTTPClient, cfg.BaseURL, authenticator),
	}

	return client, nil
}

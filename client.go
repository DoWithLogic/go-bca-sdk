package bca

import (
	"github.com/DoWithLogic/go-bca-sdk/account_information"
	"github.com/DoWithLogic/go-bca-sdk/business_debit_card"
	"github.com/DoWithLogic/go-bca-sdk/internal/transport"
)

// Client is the main client for interacting with the BCA API.
type Client struct {
	config Config

	AccountInformation *account_information.AccountInformationService
	BusinessDebitCard  *business_debit_card.BusinessDebitCardService
}

// NewClient creates a new BCA API client using the provided options.
func NewClient(opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	authenticator, err := cfg.newAuthenticator()
	if err != nil {
		return nil, err
	}

	transport := transport.NewClient(
		cfg.HTTPClient,
		cfg.BaseURL,
		authenticator,
		transport.RetryConfig{MaxRetries: cfg.MaxRetries, Backoff: cfg.RetryBackoff},
	)

	client := &Client{
		config:             cfg,
		AccountInformation: account_information.NewAccountInformationService(transport),
		BusinessDebitCard:  business_debit_card.NewBusinessDebitCardService(transport),
	}

	return client, nil
}

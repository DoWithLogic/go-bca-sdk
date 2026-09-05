package transport

import (
	"net/http"
	"time"

	"github.com/DoWithLogic/go-bca-sdk/internal/auth"
)

// Client is an HTTP client used to communicate with the BCA API.
type Client struct {
	httpClient  *http.Client
	baseURL     string
	auth        auth.Authenticator
	retryPolicy RetryPolicy
	retryConfig RetryConfig
}

// NewClient creates a new HTTP client for the BCA API.
//
// The provided HTTP client is used to execute API requests, and baseURL
// specifies the base URL of the BCA API.
func NewClient(httpClient *http.Client, baseURL string, authenticator auth.Authenticator) *Client {
	return &Client{
		httpClient:  httpClient,
		baseURL:     baseURL,
		auth:        authenticator,
		retryPolicy: DefaultRetryPolicy{},
		retryConfig: RetryConfig{
			MaxRetries: 2,
			Backoff:    100 * time.Microsecond,
		},
	}
}

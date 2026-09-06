package transport

import (
	"net/http"
	"time"
)

// RetryPolicy determines whether an HTTP request should be retried.
type RetryPolicy interface {
	ShouldRetry(method string, statusCode int) bool
}

// DefaultRetryPolicy determines which HTTP responses are safe to retry.
type DefaultRetryPolicy struct{}

// ShouldRetry return true for transient server errors.
func (DefaultRetryPolicy) ShouldRetry(method string, statusCode int) bool {
	switch method {
	case http.MethodGet, http.MethodHead:
		return statusCode == http.StatusRequestTimeout || statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
	default:
		return false
	}
}

// RetryConfig configures request retries.
type RetryConfig struct {
	MaxRetries int
	Backoff    time.Duration
}

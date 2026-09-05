package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Do executes an HTTP request against the BCA API.
//
// The body is encoded as JSON when it is not nil. If result is not nil,
// the response body is decoded from JSON into result.
func (c *Client) Do(ctx context.Context, method string, path string, body any, result any) error {
	url := c.baseURL + path

	var requestBody []byte
	var err error

	if body != nil {
		requestBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	for attempt := 0; ; attempt++ {
		req, err := http.NewRequestWithContext(
			ctx,
			method,
			url,
			bytes.NewReader(requestBody),
		)
		if err != nil {
			return err
		}

		req.Header.Set("Accept", "application/json")

		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		if c.auth != nil {
			if err := c.auth.Authenticate(ctx, req); err != nil {
				return err
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode >= http.StatusOK &&
			resp.StatusCode < http.StatusMultipleChoices {
			defer resp.Body.Close()

			if result == nil {
				return nil
			}

			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return err
			}

			return nil
		}

		responseBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return err
		}

		if !c.retryPolicy.ShouldRetry(method, resp.StatusCode) ||
			attempt >= c.retryConfig.MaxRetries {
			return &APIError{
				StatusCode: resp.StatusCode,
				Status:     resp.Status,
				Body:       responseBody,
			}
		}

		backoff := c.retryConfig.Backoff * time.Duration(1<<attempt)

		if err := sleepWithContext(ctx, backoff); err != nil {
			return err
		}
	}
}

func sleepWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

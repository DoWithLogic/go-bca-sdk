# Users API SDK for Go

Unofficial Go SDK for integrating with the BCA Developer API.

The goal of this project is to provide an idiomatic, production-ready Go client for the BCA Developer API, with support for authentication, request signing, HTTP communication, and type-safe API clients.

> **Disclaimer:** This is an unofficial, community-developed SDK and is not
> affiliated with or endorsed by Bank Central Asia (BCA).

---

## Features

- **OAuth 2.0 Bearer Token Authentication**
- **Type-Safe API Clients** - Strongly typed requests and responses
- **Context-Aware Requests** - Full support for cancellation and timeouts
- **Automatic Retry Mechanisms** - Configurable retry for transient failures
- **Comprehensive Error Handling** - Specific error types for better error handling
- **Well-Documented** - Clear examples and API reference
- **Testable** - Easy to mock and test
- **RESTful Operations** - Full CRUD support for resources

---

## Installation

```bash
go get github.com/DoWithLogic/go-bca-sdk
```

---

## Quick Start

### SNAP Authentication

```go
package main

import (
	"crypto/rsa"
	"log"

	bca "github.com/DoWithLogic/go-bca-sdk"
)

func main() {
	var privateKey *rsa.PrivateKey

	client, err := bca.NewClient(
		bca.WithClientID("your-client-id"),
		bca.WithClientSecret("your-client-secret"),
		bca.WithSNAPAuth(privateKey),
		bca.WithSNAPChannelID("95051"),
		bca.WithSNAPPartnerID("your-partner-id"),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = client
}
```

### Balance Inquiry

Perform a SNAP Banking Balance Inquiry for a registered KlikBCA Bisnis account:

```go
package main

import (
	"context"
	"crypto/rsa"
	"log"

	bca "github.com/DoWithLogic/go-bca-sdk"
)

func main() {
	var privateKey *rsa.PrivateKey

	client, err := bca.NewClient(
		bca.WithClientID("your-client-id"),
		bca.WithClientSecret("your-client-secret"),
		bca.WithSNAPAuth(privateKey),
		bca.WithSNAPChannelID("95051"),
		bca.WithSNAPPartnerID("your-partner-id"),
	)
	if err != nil {
		log.Fatal(err)
	}

	request := bca.BalanceInquiryRequest{
		PartnerReferenceNo: "partner-reference-123",
		AccountNo:          "1234567890",
	}

	response, err := client.Account.BalanceInquiry(
		context.Background(),
		request,
		"external-reference-123",
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(response.ResponseCode)
	log.Println(response.ResponseMessage)
	log.Println(response.AccountNo)
	log.Println(response.Name)
}
```

`X-EXTERNAL-ID` is provided by the consumer and should be unique for the same day.

## Error Handling

API errors are returned as `*APIError` and contain both the HTTP response information and the BCA response information.

### Checking Specific Error Types

```go
response, err := client.Account.BalanceInquiry(
	context.Background(),
	request,
	"external-reference-123",
)
if err != nil {
	var apiErr *bca.APIError
	if errors.As(err, &apiErr) {
		log.Println(apiErr.StatusCode)
		log.Println(apiErr.ResponseCode)
		log.Println(apiErr.ResponseMessage)
		log.Println(string(apiErr.Body))
	}

	return
}
```

The raw response body is preserved in `APIError.Body`, allowing consumers to access the original BCA error response.

---

## Retry Behavior

The SDK supports configurable retries for retryable HTTP requests.

By default, retries are enabled for:

- `408 Request Timeout`
- `429 Too Many Requests`
- `5xx Server Errors`

Retries are limited to idempotent HTTP methods by default. POST requests are not automatically retried to avoid unintentionally repeating operations such as money movement.

## Status

🚧 **Work in Progress**

The SDK is currently under active development. APIs and implementation details may change before the first stable release.

---

## Documentation

- [BCA Developer API Documentation](https://developer.bca.co.id/Dokumentasi)

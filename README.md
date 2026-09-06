# BCA Developer API SDK for Go

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

Complete examples are available in the [`examples/`](examples/) directory:

- [`examples/oauth2`](examples/oauth2) — BCA OAuth 2.0 authentication
- [`examples/snap`](examples/snap) — BCA SNAP authentication

## Error Handling

API errors are returned as `*APIError` and contain both the HTTP response information and the BCA response information.

### Checking Specific Error Types

```go
response, err := client.AccountInformation.BalanceInquiry(
    context.Background(),
    request,
    "external-reference-123",
)
if err != nil {
    var apiErr *bcaerrors.APIError
    if errors.As(err, &apiErr) {
        log.Println(apiErr.HTTPStatusCode)
        log.Println(apiErr.ResponseCode)
        log.Println(apiErr.ResponseMessage)
    }

    return
}
```

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

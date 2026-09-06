# Contributing

Thank you for your interest in contributing to the BCA Developer API SDK for Go.

Contributions, bug reports, documentation improvements, and feature suggestions are welcome.

## Getting Started

### 1. Fork the Repository

Create a fork of the repository and clone it locally:

```bash
git clone https://github.com/DoWithLogic/go-bca-sdk.git
cd go-bca-sdk
```

### 2. Create a Branch

Create a branch for your changes:

```bash
git checkout -b feature/your-feature
```

Use a descriptive branch name, for example:

- `feature/snap-balance-inquiry`
- `fix/oauth-token-refresh`
- `docs/update-readme`

## Development

Make sure your changes are formatted and tested before submitting a pull request.

Run:

```bash
go fmt ./...
go vet ./...
go test ./...
```

If you add or modify functionality, please include appropriate tests.

## Adding a New API

When adding support for a new BCA API:

1. Follow the API structure defined in the BCA Developer API documentation.
2. Add strongly typed request and response models.
3. Add the corresponding service method.
4. Add tests for successful and error responses.
5. Update the documentation or examples when appropriate.

Keep public APIs simple and consistent with the existing SDK design.

## Error Handling

Use the existing `errors.APIError` for BCA API errors rather than introducing a new error type for individual endpoints.

For example:

```go
var apiErr *bcaerrors.APIError

if errors.As(err, &apiErr) {
    // Handle BCA API error
}
```

## Pull Requests

Before opening a pull request:

- Ensure all tests pass.
- Run `go fmt ./...`.
- Run `go vet ./...`.
- Add tests for new functionality.
- Keep changes focused and related to the purpose of the pull request.
- Update documentation when necessary.

Please provide a clear description of:

- What changed
- Why the change was needed
- How it was tested

## Commit Messages

Use clear and descriptive commit messages.

Examples:

```text
feat: add SNAP balance inquiry
fix: handle expired OAuth access token
test: add debit card inquiry error cases
docs: update authentication examples
```

## Code Style

Follow standard Go conventions and prefer simple, readable implementations.

When in doubt, follow the conventions already established in the repository.

## Reporting Issues

When reporting a bug, please include:

- Go version
- SDK version
- Relevant API/endpoint
- Expected behavior
- Actual behavior
- A minimal reproduction when possible

**Do not include credentials, API secrets, private keys, access tokens, or other sensitive information in issues or pull requests.**

package virtual_account

import "github.com/DoWithLogic/go-bca-sdk/internal/transport"

// VirtualAccountService provides access to the Virtual Account APIs.
type VirtualAccountService struct {
	t *transport.Client
}

// NewVirtualAccountService creates a new VirtualAccountService
// using the provided transport client.
func NewVirtualAccountService(t *transport.Client) *VirtualAccountService {
	return &VirtualAccountService{t: t}
}

package account_information

import "github.com/DoWithLogic/go-bca-sdk/internal/transport"

// AccountInformationService provides account-related services through the BCA API.
type AccountInformationService struct {
	transport *transport.Client
}

func NewAccountInformationService(t *transport.Client) *AccountInformationService {
	return &AccountInformationService{t}
}

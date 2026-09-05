package business_debit_card

import "github.com/DoWithLogic/go-bca-sdk/internal/transport"

type BusinessDebitCardService struct {
	transport *transport.Client
}

func NewBusinessDebitCardService(t *transport.Client) *BusinessDebitCardService {
	return &BusinessDebitCardService{t}
}

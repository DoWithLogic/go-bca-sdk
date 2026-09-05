package bca

import (
	"context"
	"net/http"

	"github.com/DoWithLogic/go-bca-sdk/internal/transport"
)

// AccountService provides account-related services through the BCA API.
type AccountService struct {
	client *Client
}

// BalanceInquiryRequest contains the parameters required to inquire an
// account balance through the BCA SNAP Banking Balance Inquiry service.
type BalanceInquiryRequest struct {
	PartnerReferenceNo string `json:"partnerReferenceNo"`
	AccountNo          string `json:"accountNo"`
}

// BalanceInquiryResponse represents the response from the BCA SNAP Banking
// Balance Inquiry service.
type BalanceInquiryResponse struct {
	ResponseCode       string        `json:"responseCode"`
	ResponseMessage    string        `json:"responseMessage"`
	ReferenceNo        string        `json:"referenceNo"`
	PartnerReferenceNo string        `json:"partnerReferenceNo"`
	AccountNo          string        `json:"accountNo"`
	Name               string        `json:"name"`
	AccountInfos       []AccountInfo `json:"accountInfos"`
}

// AccountInfo contains balance information for an account.
type AccountInfo struct {
	Amount           *Money `json:"amount"`
	FloatAmount      *Money `json:"floatAmount"`
	HoldAmount       *Money `json:"holdAmount"`
	AvailableBalance *Money `json:"availableBalance"`
}

// Money represents a monetary amount and its ISO-4217 currency.
type Money struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// BalanceInquiry performs a SNAP Banking Balance Inquiry for a registered
// KlikBCA Bisnis account.
func (s *AccountService) BalanceInquiry(ctx context.Context, request BalanceInquiryRequest, externalID string) (*BalanceInquiryResponse, error) {
	var response BalanceInquiryResponse
	err := s.client.transport.Do(
		ctx,
		transport.Request{
			Method: http.MethodPost,
			Path:   "/openapi/v1.0/balance-inquiry",
			Headers: http.Header{
				"X-EXTERNAL-ID": []string{externalID},
			},
			Body: request,
		},
		&response,
	)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

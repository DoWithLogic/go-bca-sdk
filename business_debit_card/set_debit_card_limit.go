package business_debit_card

import (
	"context"
	"net/http"

	"github.com/DoWithLogic/go-bca-sdk/internal/transport"
)

// SetDebitCardLimitRequest represents a request to set the limit
// for a Business Debit Card (BDC).
type SetDebitCardLimitRequest struct {
	// TransactionID is the transaction identifier.
	TransactionID string `json:"transactionId"`

	// EffectiveDate is the effective date of the limit setting in YYYYMMDD format.
	EffectiveDate string `json:"effectiveDate"`

	// AccountNumber is the BDC account number.
	AccountNumber string `json:"accountNumber"`

	// CardType specifies the BDC type: L (Loyalty), P (Petty Cash), or D (Deposit).
	CardType string `json:"cardType"`

	// RecurringType specifies how the limit is reset:
	// N (No Reset), D (Reset by Date), or X (Reset by Day).
	RecurringType string `json:"recurringType"`

	// RecurringPeriod specifies the recurring period for the BDC limit.
	RecurringPeriod string `json:"recurringPeriod"`

	// CardNumber is the 16-digit BDC card number.
	CardNumber string `json:"cardNumber"`

	// ExpiryLimit is the expiration date of the limit in YYYYMMDD format.
	ExpiryLimit string `json:"expiryLimit"`

	// LimitAmount specifies the BDC limit amount.
	// A numeric value changes the limit, while a space leaves the existing limit unchanged.
	LimitAmount string `json:"limitAmount"`

	// UniqueKey is the customer's unique key.
	UniqueKey string `json:"uniqueKey"`

	// Cardholder is the name of the BDC cardholder.
	Cardholder string `json:"cardholder,omitempty"`

	// IDNumber is the identification number of the BDC cardholder.
	IDNumber string `json:"idNumber,omitempty"`

	// PhoneNumber is the phone number of the BDC cardholder.
	PhoneNumber string `json:"phoneNumber,omitempty"`

	// EmailAddress is the email address of the BDC cardholder.
	EmailAddress string `json:"emailAddress,omitempty"`

	// Description contains additional information about the limit setting.
	Description string `json:"description,omitempty"`

	// ProcessType specifies the limit setting process type.
	// Only S (Single Process) is supported.
	ProcessType string `json:"processType"`
}

// SetDebitCardLimitResponse represents the response from the
// Business Debit Card Set Limit API.
type SetDebitCardLimitResponse struct {
	// ResponseCode identifies the transaction status.
	ResponseCode string `json:"responseCode"`

	// ResponseMessage contains the description of the transaction result.
	ResponseMessage string `json:"responseMessage"`

	// ResponseData contains the transaction result details.
	ResponseData SetDebitCardLimitResponseData `json:"responseData"`
}

// SetDebitCardLimitResponseData contains the result details
// returned by the Business Debit Card Set Limit API.
type SetDebitCardLimitResponseData struct {
	// TransactionID is the transaction identifier assigned by the service provider.
	TransactionID string `json:"transactionId"`

	// UpdateDateTime is the transaction date and time in YYYYMMDDHHMMSS format.
	UpdateDateTime string `json:"updateDateTime"`

	// CardNumber is the BDC card number.
	CardNumber string `json:"cardNumber"`
}

// SetDebitCardLimit sets or updates the limit for a Business Debit Card (BDC).
//
// The externalID must be unique for the same day.
func (bdc *BusinessDebitCardService) SetDebitCardLimit(ctx context.Context, request SetDebitCardLimitRequest, externalID string) (*SetDebitCardLimitResponse, error) {
	var response SetDebitCardLimitResponse
	err := bdc.transport.Do(
		ctx,
		transport.Request{
			Method: http.MethodPost,
			Path:   "/openapi/bdc/v1.0/setting-limit",
			Body:   request,
			Headers: http.Header{
				"X-EXTERNAL-ID": []string{externalID},
			},
		},
		&response,
	)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

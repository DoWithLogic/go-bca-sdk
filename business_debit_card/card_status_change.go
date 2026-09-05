package business_debit_card

import (
	"context"
	"net/http"

	"github.com/DoWithLogic/go-bca-sdk/internal/transport"
)

// BlockType specifies the block status to apply to a BDC card.
type BlockType string

const (
	// BlockTypeOpenHalf indicates an open half block.
	BlockTypeOpenHalf BlockType = "O"

	// BlockTypeHalf indicates a half block.
	BlockTypeHalf BlockType = "H"

	// BlockTypeFull indicates a full block.
	BlockTypeFull BlockType = "B"
)

// ReasonCode specifies the reason for changing the BDC card status.
type ReasonCode string

const (
	// ReasonCodeLost indicates that the card was lost.
	ReasonCodeLost ReasonCode = "H"

	// ReasonCodeSwallowed indicates that the card was swallowed.
	ReasonCodeSwallowed ReasonCode = "R"

	// ReasonCodeOthers indicates other reasons.
	ReasonCodeOthers ReasonCode = "O"
)

// CardStatusChangeRequest represents a request to change the status of a BDC card.
type CardStatusChangeRequest struct {
	// CardNumber is the 16-digit BDC card number.
	CardNumber string `json:"cardNumber"`

	// BlockType specifies the block status to apply to the card.
	BlockType BlockType `json:"blockType"`

	// ReasonCode specifies the reason for changing the card status.
	ReasonCode ReasonCode `json:"reasonCode"`
}

// CardStatusChangeResponse represents the response from a BDC card status change request.
type CardStatusChangeResponse struct {
	// ResponseCode identifies the transaction status.
	//
	// The format is AAABBCC:
	// AAA: HTTP status code
	// BB: Service code
	// CC: Case code
	//
	// The service code for BDC Update Status is US.
	ResponseCode string `json:"responseCode"`

	// ResponseMessage describes the transaction result.
	ResponseMessage string `json:"responseMessage"`

	// ResponseData contains the updated BDC card information.
	ResponseData CardStatusChangeResponseData `json:"responseData"`
}

// CardStatusChangeResponseData contains the BDC card information returned
// by a card status change request.
type CardStatusChangeResponseData struct {
	// CardNumber is the 16-digit BDC card number.
	//
	// If the request fails, this field is returned as "-".
	CardNumber string `json:"cardNumber"`

	// CardStatus is the current status of the BDC card.
	//
	// If the request fails, this field is returned as "-".
	CardStatus string `json:"cardStatus"`
}

// CardStatusChange updates the status of a BDC card.
func (bdc *BusinessDebitCardService) CardStatusChange(ctx context.Context, request CardStatusChangeRequest, externalID string) (*CardStatusChangeResponse, error) {
	var response CardStatusChangeResponse

	err := bdc.transport.Do(
		ctx,
		transport.Request{
			Method:  http.MethodPost,
			Path:    "/openapi/bdc/v1.0/update-status",
			Body:    request,
			Headers: http.Header{"X-EXTERNAL-ID": []string{externalID}},
		},
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

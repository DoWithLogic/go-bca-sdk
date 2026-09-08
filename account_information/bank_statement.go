package account_information

import (
	"context"
	"net/http"

	"github.com/DoWithLogic/go-bca-sdk/internal/transport"
)

// BankStatementRequest represents a request to retrieve
// account transaction statements for a specified time range.
type BankStatementRequest struct {
	PartnerReferenceNo string `json:"partnerReferenceNo"`
	AccountNo          string `json:"accountNo"`
	FromDateTime       string `json:"fromDateTime"`
	ToDateTime         string `json:"toDateTime"`
}

// BankStatementResponse represents the response from the
// SNAP Banking Bank Statement API.
type BankStatementResponse struct {
	ResponseCode       string                `json:"responseCode"`
	ResponseMessage    string                `json:"responseMessage"`
	ReferenceNo        string                `json:"referenceNo"`
	PartnerReferenceNo string                `json:"partnerReferenceNo"`
	Balance            []Balance             `json:"balance"`
	TotalCreditEntries StatementEntries      `json:"totalCreditEntries"`
	TotalDebitEntries  StatementEntries      `json:"totalDebitEntries"`
	DetailData         []BankStatementDetail `json:"detailData"`
}

// Balance represents an account balance at a specific point in time.
type Balance struct {
	Amount          Amount `json:"amount"`
	StartingBalance Amount `json:"startingBalance"`
	EndingBalance   Amount `json:"endingBalance"`
}

// Amount represents a monetary amount and its currency.
type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// StatementEntries represents the total number and amount
// of credit or debit transactions.
type StatementEntries struct {
	NumberOfEntries string `json:"numberOfEntries"`
	Amount          Amount `json:"amount"`
}

// BankStatementDetail represents an individual account transaction.
type BankStatementDetail struct {
	Amount          Amount `json:"amount"`
	TransactionDate string `json:"transactionDate"`
	Remark          string `json:"remark"`
	Type            string `json:"type"`
}

// BankStatement retrieves account transaction statements
// for the specified time range.
func (ai *AccountInformationService) BankStatement(ctx context.Context, request BankStatementRequest, externalID string) (*BankStatementResponse, error) {
	var response BankStatementResponse
	err := ai.transport.Do(
		ctx,
		transport.Request{
			Method: http.MethodPost,
			Path:   "/openapi/v1.0/bank-statement",
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

package business_debit_card

import (
	"context"
	"net/http"

	"github.com/DoWithLogic/go-bca-sdk/internal/transport"
)

// DebitCardInquiryRequest contains the parameters required to inquire
// about a BDC debit card.
type DebitCardInquiryRequest struct {
	BankCardToken string `json:"bankCardToken"`
}

// DebitCardInquiryResponse represents the response from the BCA BDC
// Debit Card Inquiry service.
type DebitCardInquiryResponse struct {
	ResponseCode    string               `json:"responseCode"`
	ResponseMessage string               `json:"responseMessage"`
	ResponseData    DebitCardInquiryData `json:"responseData"`
}

// DebitCardInquiryData contains the BDC debit card information.
type DebitCardInquiryData struct {
	AccountNo      string                   `json:"accountNo"`
	Name           string                   `json:"name"`
	AccountInfos   []DebitCardAccountInfo   `json:"accountInfos"`
	AdditionalInfo *DebitCardAdditionalInfo `json:"additionalInfo"`
}

// DebitCardAccountInfo contains BDC account balance information.
type DebitCardAccountInfo struct {
	BalanceType string `json:"balanceType"`
	Amount      Money  `json:"amount"`
}

// DebitCardAdditionalInfo contains additional BDC card information.
type DebitCardAdditionalInfo struct {
	CardNumber      string `json:"cardNumber"`
	FactType        string `json:"factType"`
	TotalAccount    string `json:"totalAccount"`
	LanID           string `json:"lanId"`
	OpenDate        string `json:"openDate"`
	ProdDate        string `json:"prodDate"`
	CloseDate       string `json:"closeDate"`
	ExpiryDate      string `json:"expiryDate"`
	EmbName         string `json:"embName"`
	Addr1           string `json:"addr1"`
	Addr2           string `json:"addr2"`
	Addr3           string `json:"addr3"`
	BirthDate       string `json:"birthDate"`
	BranchCode      string `json:"branchCode"`
	EmpCode         string `json:"empCode"`
	ChargeFlag      string `json:"chargeFlag"`
	Type            string `json:"type"`
	DeptCode        string `json:"deptCode"`
	StockBranchCode string `json:"stockBranchCode"`
	CustomerCode    string `json:"customerCode"`
	CustomField     string `json:"customField"`
	CardStatus      string `json:"cardStatus"`
	CorpID          string `json:"corpId"`
	Outstanding     string `json:"outstanding"`
	SendFlag        string `json:"sendFlag"`
	IDNumber        string `json:"idNumber"`
	CustomerID      string `json:"customerId"`
	MaidenName      string `json:"maidenName"`
	CorpType        string `json:"corpType"`
	UniqueKey       string `json:"uniqueKey"`
	LimitUsage      string `json:"limitUsage"`
	LimitNextDate   string `json:"limitNextDate"`
	LimitExpiryDate string `json:"limitExpiryDate"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Desc            string `json:"desc"`
	LogonID         string `json:"logonId"`
	ReasonCode      string `json:"reasonCode"`
	ResetType       string `json:"resetType"`
	ResetPeriod     string `json:"resetPeriod"`
	PhoneEcomm      string `json:"phoneEcomm"`
}

// Money represents a monetary amount and its ISO 4217 currency.
type Money struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// DebitCardInquiry performs a BDC Debit Card Inquiry request.
func (bdc *BusinessDebitCardService) DebitCardInquiry(ctx context.Context, request DebitCardInquiryRequest, externalID string) (*DebitCardInquiryResponse, error) {
	var response DebitCardInquiryResponse
	err := bdc.transport.Do(
		ctx,
		transport.Request{
			Method:  http.MethodPost,
			Path:    "/openapi/bdc/v1.0/inquiry-card",
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

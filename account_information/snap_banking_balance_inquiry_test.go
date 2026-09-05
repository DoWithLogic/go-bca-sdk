package account_information

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DoWithLogic/go-bca-sdk/internal/transport"
)

type mockAuthenticator struct{}

func (mockAuthenticator) Authenticate(
	ctx context.Context,
	req *http.Request,
) error {
	return nil
}

func TestAccountInformationService_BalanceInquiry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method %s, got %s", http.MethodPost, r.Method)
		}

		if r.URL.Path != "/openapi/v1.0/balance-inquiry" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if got := r.Header.Get("X-EXTERNAL-ID"); got != "external-123" {
			t.Errorf("expected X-EXTERNAL-ID %q, got %q", "external-123", got)
		}

		var request BalanceInquiryRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		expectedRequest := BalanceInquiryRequest{
			PartnerReferenceNo: "partner-123",
			AccountNo:          "1234567890",
		}

		if request != expectedRequest {
			t.Errorf("unexpected request: %+v", request)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(BalanceInquiryResponse{
			ResponseCode:       "2001100",
			ResponseMessage:    "Successful",
			ReferenceNo:        "2609050100000001",
			PartnerReferenceNo: "partner-123",
			AccountNo:          "1234567890",
			Name:               "MARTIN",
			AccountInfos: []AccountInfo{
				{
					Amount: &Money{
						Value:    "1000000.00",
						Currency: "IDR",
					},
					FloatAmount: &Money{
						Value:    "0.00",
						Currency: "IDR",
					},
					HoldAmount: &Money{
						Value:    "100000.00",
						Currency: "IDR",
					},
					AvailableBalance: &Money{
						Value:    "900000.00",
						Currency: "IDR",
					},
				},
			},
		})
	}))
	defer server.Close()

	transportClient := transport.NewClient(server.Client(), server.URL, mockAuthenticator{}, transport.RetryConfig{MaxRetries: 0})

	service := &AccountInformationService{
		transport: transportClient,
	}

	request := BalanceInquiryRequest{
		PartnerReferenceNo: "partner-123",
		AccountNo:          "1234567890",
	}

	response, err := service.BalanceInquiry(context.Background(), request, "external-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.ResponseCode != "2001100" {
		t.Errorf("expected response code %q, got %q", "2001100", response.ResponseCode)
	}

	if response.ResponseMessage != "Successful" {
		t.Errorf("expected response message %q, got %q", "Successful", response.ResponseMessage)
	}

	if response.ReferenceNo != "2609050100000001" {
		t.Errorf("expected reference number %q, got %q", "2609050100000001", response.ReferenceNo)
	}

	if response.PartnerReferenceNo != "partner-123" {
		t.Errorf("expected partner reference number %q, got %q", "partner-123", response.PartnerReferenceNo)
	}

	if response.AccountNo != "1234567890" {
		t.Errorf("expected account number %q, got %q", "1234567890", response.AccountNo)
	}

	if response.Name != "MARTIN" {
		t.Errorf("expected name %q, got %q", "MARTIN", response.Name)
	}

	if len(response.AccountInfos) != 1 {
		t.Fatalf("expected 1 account info, got %d", len(response.AccountInfos))
	}

	accountInfo := response.AccountInfos[0]

	if accountInfo.Amount == nil {
		t.Fatal("expected amount, got nil")
	}

	if accountInfo.Amount.Value != "1000000.00" {
		t.Errorf("expected amount value %q, got %q", "1000000.00", accountInfo.Amount.Value)
	}

	if accountInfo.Amount.Currency != "IDR" {
		t.Errorf("expected amount currency %q, got %q", "IDR", accountInfo.Amount.Currency)
	}

	if accountInfo.AvailableBalance == nil {
		t.Fatal("expected available balance, got nil")
	}

	if accountInfo.AvailableBalance.Value != "900000.00" {
		t.Errorf("expected available balance %q, got %q", "900000.00", accountInfo.AvailableBalance.Value)
	}
}

func TestAccountInformationService_BalanceInquiry_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"responseCode":    "4091100",
			"responseMessage": "Conflict",
		})
	}))
	defer server.Close()

	transportClient := transport.NewClient(server.Client(), server.URL, mockAuthenticator{}, transport.RetryConfig{MaxRetries: 0})
	service := &AccountInformationService{transport: transportClient}

	_, err := service.BalanceInquiry(
		context.Background(),
		BalanceInquiryRequest{
			PartnerReferenceNo: "partner-123",
			AccountNo:          "1234567890",
		},
		"external-123",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

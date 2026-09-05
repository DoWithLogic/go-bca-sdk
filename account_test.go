package bca

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DoWithLogic/go-bca-sdk/internal/transport"
)

func TestAccountService_BalanceInquiry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/openapi/v1.0/balance-inquiry" {
			t.Fatalf("expected /openapi/v1.0/balance-inquiry, got %s", r.URL.Path)
		}

		var request BalanceInquiryRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if request.PartnerReferenceNo != "202609050001" {
			t.Fatalf("unexpected partner reference number: %s", request.PartnerReferenceNo)
		}

		if request.AccountNo != "1234567890" {
			t.Fatalf("unexpected account number: %s", request.AccountNo)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"responseCode": "2001100",
			"responseMessage": "Successful",
			"referenceNo": "20260905000101",
			"partnerReferenceNo": "202609050001",
			"accountNo": "1234567890",
			"name": "JOHN DOE",
			"accountInfos": [
				{
					"amount": {
						"value": "10000000.00",
						"currency": "IDR"
					},
					"floatAmount": {
						"value": "0.00",
						"currency": "IDR"
					},
					"holdAmount": {
						"value": "100000.00",
						"currency": "IDR"
					},
					"availableBalance": {
						"value": "9900000.00",
						"currency": "IDR"
					}
				}
			]
		}`))
	}))
	defer server.Close()

	cfg := defaultConfig()
	cfg.BaseURL = server.URL
	cfg.AuthMode = AuthModeSNAP

	client := &Client{
		config:    cfg,
		transport: transport.NewClient(cfg.HTTPClient, cfg.BaseURL, nil, transport.RetryConfig{}),
	}

	client.Account = &AccountService{
		client: client,
	}

	response, err := client.Account.BalanceInquiry(context.Background(), BalanceInquiryRequest{PartnerReferenceNo: "202609050001", AccountNo: "1234567890"}, "28910000006578499987546738976812")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.ResponseCode != "2001100" {
		t.Fatalf("unexpected response code: %s", response.ResponseCode)
	}

	if response.ResponseMessage != "Successful" {
		t.Fatalf("unexpected response message: %s", response.ResponseMessage)
	}

	if response.ReferenceNo != "20260905000101" {
		t.Fatalf("unexpected reference number: %s", response.ReferenceNo)
	}

	if response.Name != "JOHN DOE" {
		t.Fatalf("unexpected account name: %s", response.Name)
	}

	if len(response.AccountInfos) != 1 {
		t.Fatalf("expected 1 account info, got %d", len(response.AccountInfos))
	}

	accountInfo := response.AccountInfos[0]
	if accountInfo.AvailableBalance.Value != "9900000.00" {
		t.Fatalf("unexpected available balance: %s", accountInfo.AvailableBalance.Value)
	}

	if accountInfo.AvailableBalance.Currency != "IDR" {
		t.Fatalf("unexpected currency: %s", accountInfo.AvailableBalance.Currency)
	}
}

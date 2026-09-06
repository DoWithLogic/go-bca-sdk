package business_debit_card

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	bcaErrors "github.com/DoWithLogic/go-bca-sdk/errors"

	"github.com/DoWithLogic/go-bca-sdk/internal/transport"
)

func TestBusinessDebitCardService_DebitCardInquiry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/openapi/bdc/v1.0/inquiry-card" {
			t.Errorf("expected path %q, got %q", "/openapi/bdc/v1.0/inquiry-card", r.URL.Path)
		}

		if got := r.Header.Get("X-EXTERNAL-ID"); got != "01234506022024001" {
			t.Errorf("expected X-EXTERNAL-ID %q, got %q", "01234506022024001", got)
		}

		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", got)
		}

		var request DebitCardInquiryRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if request.BankCardToken != "1234567890123456" {
			t.Errorf("expected bankCardToken %q, got %q", "1234567890123456", request.BankCardToken)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{
			"responseCode": "200IC00",
			"responseMessage": "Successful",
			"responseData": {
				"accountNo": "",
				"name": "John Doe",
				"accountInfos": [
					{
						"balanceType": "Limit",
						"amount": {
							"value": "20000.00",
							"currency": "IDR"
						}
					}
				],
				"additionalInfo": {
					"cardNumber": "1234567890123456",
					"facType": "1",
					"totalAccount": "01",
					"lanId": "1",
					"openDate": "20080206",
					"prodDate": "20080206",
					"closeDate": "00000000",
					"expiryDate": "4912",
					"embName": "st Code Corp Regresi 5",
					"addr1": "JKT",
					"addr2": "JKT",
					"addr3": "THAMRIN",
					"birthDate": "00000000",
					"branchCode": "0437",
					"empCode": "F",
					"chargeFlag": "P",
					"type": "0",
					"deptCode": "00000",
					"stockBranchCode": "0437",
					"customerCode": "T",
					"customField": "REG CORP ENG 13",
					"cardStatus": "O",
					"corpId": "000097",
					"outstanding": "0",
					"sendFlag": "B",
					"idNumber": "2100001000001103",
					"customerId": "",
					"maidenName": "",
					"corpType": "L",
					"uniqueKey": "R000052",
					"limitUsage": "0",
					"limitNextDate": "20991231",
					"limitExpDate": "20180110",
					"phone": "08561589203",
					"email": "REG-CORP13@regresi.COM",
					"desc": "REG CORP LOYALITY 13",
					"logonId": "userkbb",
					"reasonCode": "",
					"resetType": "N",
					"resetPeriod": "000",
					"phoneEcomm": ""
				}
			}
		}`))
	}))
	defer server.Close()

	httpClient := server.Client()
	client := transport.NewClient(httpClient, server.URL, mockAuthenticator{}, transport.RetryConfig{MaxRetries: 0})
	service := NewBusinessDebitCardService(client)
	response, err := service.DebitCardInquiry(context.Background(), DebitCardInquiryRequest{BankCardToken: "1234567890123456"}, "01234506022024001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	if response.ResponseCode != "200IC00" {
		t.Errorf("expected response code %q, got %q", "200IC00", response.ResponseCode)
	}

	if response.ResponseMessage != "Successful" {
		t.Errorf("expected response message %q, got %q", "Successful", response.ResponseMessage)
	}

	if response.ResponseData.Name != "John Doe" {
		t.Errorf("expected name %q, got %q", "John Doe", response.ResponseData.Name)
	}

	if len(response.ResponseData.AccountInfos) != 1 {
		t.Fatalf("expected 1 account info, got %d", len(response.ResponseData.AccountInfos))
	}

	accountInfo := response.ResponseData.AccountInfos[0]

	if accountInfo.BalanceType != "Limit" {
		t.Errorf("expected balance type %q, got %q", "Limit", accountInfo.BalanceType)
	}

	if accountInfo.Amount.Value != "20000.00" {
		t.Errorf("expected amount %q, got %q", "20000.00", accountInfo.Amount.Value)
	}

	if accountInfo.Amount.Currency != "IDR" {
		t.Errorf("expected currency %q, got %q", "IDR", accountInfo.Amount.Currency)
	}

	if response.ResponseData.AdditionalInfo.CardNumber != "1234567890123456" {
		t.Errorf("expected card number %q, got %q", "1234567890123456", response.ResponseData.AdditionalInfo.CardNumber)
	}

	if response.ResponseData.AdditionalInfo.CardStatus != "O" {
		t.Errorf("expected card status %q, got %q", "O", response.ResponseData.AdditionalInfo.CardStatus)
	}
}

func TestBusinessDebitCardService_DebitCardInquiry_Error(t *testing.T) {
	tests := []struct {
		name         string
		httpStatus   int
		responseCode string
		responseMsg  string
	}{
		{
			name:         "bad request",
			httpStatus:   http.StatusBadRequest,
			responseCode: "400IC00",
			responseMsg:  "Bad Request",
		},
		{
			name:         "invalid field format",
			httpStatus:   http.StatusBadRequest,
			responseCode: "400IC01",
			responseMsg:  "Invalid Field Format bankCardToken",
		},
		{
			name:         "unauthorized",
			httpStatus:   http.StatusUnauthorized,
			responseCode: "401IC00",
			responseMsg:  "Unauthorized",
		},
		{
			name:         "invalid token",
			httpStatus:   http.StatusUnauthorized,
			responseCode: "401IC01",
			responseMsg:  "Invalid token (B2B)",
		},
		{
			name:         "feature not allowed",
			httpStatus:   http.StatusForbidden,
			responseCode: "40ICL01",
			responseMsg:  "Feature Not Allowed",
		},
		{
			name:         "invalid card",
			httpStatus:   http.StatusNotFound,
			responseCode: "404IC11",
			responseMsg:  "Invalid Card",
		},
		{
			name:         "conflict",
			httpStatus:   http.StatusConflict,
			responseCode: "409IC00",
			responseMsg:  "Conflict",
		},
		{
			name:         "internal server error",
			httpStatus:   http.StatusInternalServerError,
			responseCode: "500IC01",
			responseMsg:  "Internal Server Error",
		},
		{
			name:         "timeout",
			httpStatus:   http.StatusGatewayTimeout,
			responseCode: "504IC00",
			responseMsg:  "Timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.httpStatus)

				_, _ = w.Write([]byte(fmt.Sprintf(`{
					"responseCode": %q,
					"responseMessage": %q
				}`, tt.responseCode, tt.responseMsg)))
			}))
			defer server.Close()

			httpClient := server.Client()
			client := transport.NewClient(httpClient, server.URL, mockAuthenticator{}, transport.RetryConfig{MaxRetries: 0})
			service := NewBusinessDebitCardService(client)

			response, err := service.DebitCardInquiry(context.Background(), DebitCardInquiryRequest{BankCardToken: "1234567890123456"}, "01234506022024001")
			if response != nil {
				t.Fatalf("expected nil response, got %+v", response)
			}

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var apiErr *bcaErrors.APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected *APIError, got %T", err)
			}

			if apiErr.HTTPStatusCode != tt.httpStatus {
				t.Errorf("expected HTTP status %d, got %d", tt.httpStatus, apiErr.HTTPStatusCode)
			}

			if apiErr.ResponseCode != tt.responseCode {
				t.Errorf("expected response code %q, got %q", tt.responseCode, apiErr.ResponseCode)
			}

			if apiErr.ResponseMessage != tt.responseMsg {
				t.Errorf("expected response message %q, got %q", tt.responseMsg, apiErr.ResponseMessage)
			}
		})
	}
}

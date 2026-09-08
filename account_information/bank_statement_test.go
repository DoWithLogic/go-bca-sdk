package account_information

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

func TestAccountInformationService_BankStatement(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/openapi/v1.0/bank-statement" {
			t.Errorf("expected path /openapi/v1.0/bank-statement, got %s", r.URL.Path)
		}

		if got := r.Header.Get("X-EXTERNAL-ID"); got != "28910000006578499987546738976812" {
			t.Errorf("expected X-EXTERNAL-ID header, got %s", got)
		}

		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", got)
		}

		var request BankStatementRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if request.PartnerReferenceNo != "2020102900000000000001" {
			t.Errorf("unexpected partnerReferenceNo: %s", request.PartnerReferenceNo)
		}

		if request.AccountNo != "1234567890" {
			t.Errorf("unexpected accountNo: %s", request.AccountNo)
		}

		if request.FromDateTime != "2021-04-21T00:00:00+07:00" {
			t.Errorf("unexpected fromDateTime: %s", request.FromDateTime)
		}

		if request.ToDateTime != "2021-04-21T00:00:00+07:00" {
			t.Errorf("unexpected toDateTime: %s", request.ToDateTime)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{
			"responseCode": "2001400",
			"responseMessage": "Successful",
			"referenceNo": "2020102977770000000009",
			"partnerReferenceNo": "2020102900000000000001",
			"balance": [
				{
					"amount": {
						"value": "500000.00",
						"currency": "IDR",
						"dateTime": "2023-02-02T00:00:00+07:00"
					},
					"startingBalance": {
						"value": "500000.00",
						"currency": "IDR",
						"dateTime": "2023-02-02T00:00:00+07:00"
					},
					"endingBalance": {
						"value": "500000.00",
						"currency": "IDR",
						"dateTime": "2023-02-02T00:00:00+07:00"
					}
				}
			],
			"totalCreditEntries": {
				"numberOfEntries": "10",
				"amount": {
					"value": "500000.00",
					"currency": "IDR"
				}
			},
			"totalDebitEntries": {
				"numberOfEntries": "10",
				"amount": {
					"value": "500000.00",
					"currency": "IDR"
				}
			},
			"detailData": [
				{
					"amount": {
						"value": "10000.00",
						"currency": "IDR"
					},
					"transactionDate": "2023-02-02T00:00:00+07:00",
					"remark": "TRSF E-BANKING DB",
					"type": "DEBIT"
				},
				{
					"amount": {
						"value": "25000.00",
						"currency": "IDR"
					},
					"transactionDate": "2023-02-01T00:00:00+07:00",
					"remark": "BA JASA E-BANKING",
					"type": "CREDIT"
				}
			]
		}`))
	}))
	defer server.Close()

	httpClient := server.Client()
	client := transport.NewClient(
		httpClient,
		server.URL,
		mockAuthenticator{},
		transport.RetryConfig{MaxRetries: 0},
	)

	service := NewAccountInformationService(client)

	response, err := service.BankStatement(
		context.Background(),
		BankStatementRequest{
			PartnerReferenceNo: "2020102900000000000001",
			AccountNo:          "1234567890",
			FromDateTime:       "2021-04-21T00:00:00+07:00",
			ToDateTime:         "2021-04-21T00:00:00+07:00",
		},
		"28910000006578499987546738976812",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	if response.ResponseCode != "2001400" {
		t.Errorf("unexpected responseCode: %s", response.ResponseCode)
	}

	if response.ResponseMessage != "Successful" {
		t.Errorf("unexpected responseMessage: %s", response.ResponseMessage)
	}

	if response.ReferenceNo != "2020102977770000000009" {
		t.Errorf("unexpected referenceNo: %s", response.ReferenceNo)
	}

	if response.PartnerReferenceNo != "2020102900000000000001" {
		t.Errorf("unexpected partnerReferenceNo: %s", response.PartnerReferenceNo)
	}

	if len(response.Balance) != 1 {
		t.Fatalf("expected 1 balance, got %d", len(response.Balance))
	}

	balance := response.Balance[0]

	if balance.Amount.Value != "500000.00" {
		t.Errorf("unexpected balance amount value: %s", balance.Amount.Value)
	}

	if balance.Amount.Currency != "IDR" {
		t.Errorf("unexpected balance amount currency: %s", balance.Amount.Currency)
	}

	if balance.StartingBalance.Value != "500000.00" {
		t.Errorf("unexpected starting balance value: %s", balance.StartingBalance.Value)
	}

	if balance.EndingBalance.Value != "500000.00" {
		t.Errorf("unexpected ending balance value: %s", balance.EndingBalance.Value)
	}

	if response.TotalCreditEntries.NumberOfEntries != "10" {
		t.Errorf(
			"unexpected total credit entries: %s",
			response.TotalCreditEntries.NumberOfEntries,
		)
	}

	if response.TotalCreditEntries.Amount.Value != "500000.00" {
		t.Errorf(
			"unexpected total credit amount: %s",
			response.TotalCreditEntries.Amount.Value,
		)
	}

	if response.TotalDebitEntries.NumberOfEntries != "10" {
		t.Errorf(
			"unexpected total debit entries: %s",
			response.TotalDebitEntries.NumberOfEntries,
		)
	}

	if response.TotalDebitEntries.Amount.Value != "500000.00" {
		t.Errorf(
			"unexpected total debit amount: %s",
			response.TotalDebitEntries.Amount.Value,
		)
	}

	if len(response.DetailData) != 2 {
		t.Fatalf("expected 2 detail data, got %d", len(response.DetailData))
	}

	detail := response.DetailData[0]

	if detail.Amount.Value != "10000.00" {
		t.Errorf("unexpected detail amount: %s", detail.Amount.Value)
	}

	if detail.Amount.Currency != "IDR" {
		t.Errorf("unexpected detail currency: %s", detail.Amount.Currency)
	}

	if detail.TransactionDate != "2023-02-02T00:00:00+07:00" {
		t.Errorf("unexpected transaction date: %s", detail.TransactionDate)
	}

	if detail.Remark != "TRSF E-BANKING DB" {
		t.Errorf("unexpected remark: %s", detail.Remark)
	}

	if detail.Type != "DEBIT" {
		t.Errorf("unexpected type: %s", detail.Type)
	}
}

func TestAccountInformationService_BankStatement_Error(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		responseCode    string
		responseMessage string
	}{
		{
			name:            "bad request",
			statusCode:      http.StatusBadRequest,
			responseCode:    "4001400",
			responseMessage: "Bad request",
		},
		{
			name:            "invalid mandatory field",
			statusCode:      http.StatusBadRequest,
			responseCode:    "4001402",
			responseMessage: "Invalid Mandatory Field accountNo",
		},
		{
			name:            "invalid field format",
			statusCode:      http.StatusBadRequest,
			responseCode:    "4001401",
			responseMessage: "Invalid Field Format accountNo",
		},
		{
			name:            "unauthorized",
			statusCode:      http.StatusUnauthorized,
			responseCode:    "4011400",
			responseMessage: "Unauthorized. [Reason]",
		},
		{
			name:            "invalid token",
			statusCode:      http.StatusUnauthorized,
			responseCode:    "4011401",
			responseMessage: "Invalid token (B2B)",
		},
		{
			name:            "feature not allowed",
			statusCode:      http.StatusForbidden,
			responseCode:    "4031401",
			responseMessage: "Feature Not Allowed",
		},
		{
			name:            "conflict",
			statusCode:      http.StatusConflict,
			responseCode:    "4091400",
			responseMessage: "Conflict",
		},
		{
			name:            "too many requests",
			statusCode:      http.StatusTooManyRequests,
			responseCode:    "4291400",
			responseMessage: "Too Many Requests",
		},
		{
			name:            "general error",
			statusCode:      http.StatusInternalServerError,
			responseCode:    "5001400",
			responseMessage: "General Error",
		},
		{
			name:            "internal server error",
			statusCode:      http.StatusInternalServerError,
			responseCode:    "5001401",
			responseMessage: "Internal Server Error",
		},
		{
			name:            "timeout",
			statusCode:      http.StatusGatewayTimeout,
			responseCode:    "5041400",
			responseMessage: "Timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)

				_, _ = fmt.Fprintf(w, `{
					"responseCode": "%s",
					"responseMessage": "%s"
				}`, tt.responseCode, tt.responseMessage)
			}))
			defer server.Close()

			httpClient := server.Client()
			client := transport.NewClient(
				httpClient,
				server.URL,
				mockAuthenticator{},
				transport.RetryConfig{MaxRetries: 0},
			)

			service := NewAccountInformationService(client)

			response, err := service.BankStatement(
				context.Background(),
				BankStatementRequest{
					PartnerReferenceNo: "2020102900000000000001",
					AccountNo:          "1234567890",
					FromDateTime:       "2021-04-21T00:00:00+07:00",
					ToDateTime:         "2021-04-21T00:00:00+07:00",
				},
				"28910000006578499987546738976812",
			)

			if response != nil {
				t.Errorf("expected nil response, got %+v", response)
			}

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var apiErr *bcaErrors.APIError

			if !errors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T: %v", err, err)
			}

			if apiErr.HTTPStatusCode != tt.statusCode {
				t.Errorf(
					"expected HTTP status %d, got %d",
					tt.statusCode,
					apiErr.HTTPStatusCode,
				)
			}

			if apiErr.ResponseCode != tt.responseCode {
				t.Errorf(
					"expected response code %s, got %s",
					tt.responseCode,
					apiErr.ResponseCode,
				)
			}

			if apiErr.ResponseMessage != tt.responseMessage {
				t.Errorf(
					"expected response message %s, got %s",
					tt.responseMessage,
					apiErr.ResponseMessage,
				)
			}
		})
	}
}

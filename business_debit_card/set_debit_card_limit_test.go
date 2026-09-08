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

func TestBusinessDebitCardService_SetDebitCardLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/openapi/bdc/v1.0/setting-limit" {
			t.Errorf("expected path %q, got %q", "/openapi/bdc/v1.0/setting-limit", r.URL.Path)
		}

		if got := r.Header.Get("X-EXTERNAL-ID"); got != "01234506022024001" {
			t.Errorf("expected X-EXTERNAL-ID %q, got %q", "01234506022024001", got)
		}

		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", got)
		}

		var request SetDebitCardLimitRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if request.TransactionID != "12345678901234" {
			t.Errorf("expected transactionId %q, got %q", "12345678901234", request.TransactionID)
		}

		if request.EffectiveDate != "20240901" {
			t.Errorf("expected effectiveDate %q, got %q", "20240901", request.EffectiveDate)
		}

		if request.AccountNumber != "1234567890" {
			t.Errorf("expected accountNumber %q, got %q", "1234567890", request.AccountNumber)
		}

		if request.CardType != "L" {
			t.Errorf("expected cardType %q, got %q", "L", request.CardType)
		}

		if request.RecurringType != "N" {
			t.Errorf("expected recurringType %q, got %q", "N", request.RecurringType)
		}

		if request.RecurringPeriod != "000" {
			t.Errorf("expected recurringPeriod %q, got %q", "000", request.RecurringPeriod)
		}

		if request.CardNumber != "1234567890123456" {
			t.Errorf("expected cardNumber %q, got %q", "1234567890123456", request.CardNumber)
		}

		if request.ExpiryLimit != "20241231" {
			t.Errorf("expected expiryLimit %q, got %q", "20241231", request.ExpiryLimit)
		}

		if request.LimitAmount != "5000000" {
			t.Errorf("expected limitAmount %q, got %q", "5000000", request.LimitAmount)
		}

		if request.UniqueKey != "UNIQUE-KEY-001" {
			t.Errorf("expected uniqueKey %q, got %q", "UNIQUE-KEY-001", request.UniqueKey)
		}

		if request.Cardholder != "John Doe" {
			t.Errorf("expected cardholder %q, got %q", "John Doe", request.Cardholder)
		}

		if request.IDNumber != "1234567890123456" {
			t.Errorf("expected idNumber %q, got %q", "1234567890123456", request.IDNumber)
		}

		if request.PhoneNumber != "081234567890" {
			t.Errorf("expected phoneNumber %q, got %q", "081234567890", request.PhoneNumber)
		}

		if request.EmailAddress != "john@example.com" {
			t.Errorf("expected emailAddress %q, got %q", "john@example.com", request.EmailAddress)
		}

		if request.Description != "Test limit setting" {
			t.Errorf("expected description %q, got %q", "Test limit setting", request.Description)
		}

		if request.ProcessType != "S" {
			t.Errorf("expected processType %q, got %q", "S", request.ProcessType)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{
			"responseCode": "200SL00",
			"responseMessage": "Successful",
			"responseData": {
				"transactionId": "12345678901234",
				"updateDateTime": "20240901123045",
				"cardNumber": "1234567890123456"
			}
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

	service := NewBusinessDebitCardService(client)

	response, err := service.SetDebitCardLimit(
		context.Background(),
		SetDebitCardLimitRequest{
			TransactionID:   "12345678901234",
			EffectiveDate:   "20240901",
			AccountNumber:   "1234567890",
			CardType:        "L",
			RecurringType:   "N",
			RecurringPeriod: "000",
			CardNumber:      "1234567890123456",
			ExpiryLimit:     "20241231",
			LimitAmount:     "5000000",
			UniqueKey:       "UNIQUE-KEY-001",
			Cardholder:      "John Doe",
			IDNumber:        "1234567890123456",
			PhoneNumber:     "081234567890",
			EmailAddress:    "john@example.com",
			Description:     "Test limit setting",
			ProcessType:     "S",
		},
		"01234506022024001",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	if response.ResponseCode != "200SL00" {
		t.Errorf("expected response code %q, got %q", "200SL00", response.ResponseCode)
	}

	if response.ResponseMessage != "Successful" {
		t.Errorf("expected response message %q, got %q", "Successful", response.ResponseMessage)
	}

	if response.ResponseData.TransactionID != "12345678901234" {
		t.Errorf(
			"expected transaction ID %q, got %q",
			"12345678901234",
			response.ResponseData.TransactionID,
		)
	}

	if response.ResponseData.UpdateDateTime != "20240901123045" {
		t.Errorf(
			"expected update date time %q, got %q",
			"20240901123045",
			response.ResponseData.UpdateDateTime,
		)
	}

	if response.ResponseData.CardNumber != "1234567890123456" {
		t.Errorf(
			"expected card number %q, got %q",
			"1234567890123456",
			response.ResponseData.CardNumber,
		)
	}
}

func TestBusinessDebitCardService_SetDebitCardLimit_Error(t *testing.T) {
	tests := []struct {
		name         string
		httpStatus   int
		responseCode string
		responseMsg  string
	}{
		{
			name:         "bad request",
			httpStatus:   http.StatusBadRequest,
			responseCode: "400SL00",
			responseMsg:  "Bad Request",
		},
		{
			name:         "invalid field format",
			httpStatus:   http.StatusBadRequest,
			responseCode: "400SL01",
			responseMsg:  "Invalid Field Format cardNumber",
		},
		{
			name:         "unauthorized",
			httpStatus:   http.StatusUnauthorized,
			responseCode: "401SL00",
			responseMsg:  "Unauthorized",
		},
		{
			name:         "invalid token",
			httpStatus:   http.StatusUnauthorized,
			responseCode: "401SL01",
			responseMsg:  "Invalid token (B2B)",
		},
		{
			name:         "feature not allowed",
			httpStatus:   http.StatusForbidden,
			responseCode: "403SL01",
			responseMsg:  "Feature Not Allowed",
		},
		{
			name:         "exceeds transaction amount limit",
			httpStatus:   http.StatusForbidden,
			responseCode: "403SL02",
			responseMsg:  "Exceeds Transaction Amount Limit",
		},
		{
			name:         "card blocked",
			httpStatus:   http.StatusForbidden,
			responseCode: "403SL07",
			responseMsg:  "Card Blocked",
		},
		{
			name:         "card expired",
			httpStatus:   http.StatusForbidden,
			responseCode: "403SL08",
			responseMsg:  "Card Expired",
		},
		{
			name:         "transaction not permitted",
			httpStatus:   http.StatusForbidden,
			responseCode: "403SL15",
			responseMsg:  "Transaction Not Permitted. [Reason]",
		},
		{
			name:         "invalid card",
			httpStatus:   http.StatusNotFound,
			responseCode: "404SL11",
			responseMsg:  "Invalid Card",
		},
		{
			name:         "invalid customer unique key",
			httpStatus:   http.StatusNotFound,
			responseCode: "404SL11",
			responseMsg:  "Invalid Customer uniqueKey",
		},
		{
			name:         "invalid account",
			httpStatus:   http.StatusNotFound,
			responseCode: "404SL11",
			responseMsg:  "Invalid Account",
		},
		{
			name:         "conflict",
			httpStatus:   http.StatusConflict,
			responseCode: "409SL00",
			responseMsg:  "Conflict",
		},
		{
			name:         "internal server error",
			httpStatus:   http.StatusInternalServerError,
			responseCode: "500SL01",
			responseMsg:  "Internal Server Error",
		},
		{
			name:         "timeout",
			httpStatus:   http.StatusGatewayTimeout,
			responseCode: "504SL00",
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
			client := transport.NewClient(
				httpClient,
				server.URL,
				mockAuthenticator{},
				transport.RetryConfig{MaxRetries: 0},
			)

			service := NewBusinessDebitCardService(client)

			response, err := service.SetDebitCardLimit(
				context.Background(),
				SetDebitCardLimitRequest{
					TransactionID:   "12345678901234",
					EffectiveDate:   "20240901",
					AccountNumber:   "1234567890",
					CardType:        "L",
					RecurringType:   "N",
					RecurringPeriod: "000",
					CardNumber:      "1234567890123456",
					ExpiryLimit:     "20241231",
					LimitAmount:     "5000000",
					UniqueKey:       "UNIQUE-KEY-001",
					ProcessType:     "S",
				},
				"01234506022024001",
			)

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
				t.Errorf(
					"expected HTTP status %d, got %d",
					tt.httpStatus,
					apiErr.HTTPStatusCode,
				)
			}

			if apiErr.ResponseCode != tt.responseCode {
				t.Errorf(
					"expected response code %q, got %q",
					tt.responseCode,
					apiErr.ResponseCode,
				)
			}

			if apiErr.ResponseMessage != tt.responseMsg {
				t.Errorf(
					"expected response message %q, got %q",
					tt.responseMsg,
					apiErr.ResponseMessage,
				)
			}
		})
	}
}

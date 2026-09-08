package virtual_account

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

type mockAuthenticator struct{}

func (mockAuthenticator) Authenticate(ctx context.Context, req *http.Request) error { return nil }

func TestVirtualAccountService_VirtualAccountInquiry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/openapi/v2.0/transfer-va/status" {
			t.Errorf(
				"expected path /openapi/v2.0/transfer-va/status, got %s",
				r.URL.Path,
			)
		}

		if got := r.Header.Get("X-EXTERNAL-ID"); got != "28910000006578499987546738976812" {
			t.Errorf("expected X-EXTERNAL-ID header, got %s", got)
		}

		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", got)
		}

		var request VirtualAccountInquiryRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if request.PartnerServiceID != "12345678" {
			t.Errorf("unexpected partnerServiceId: %s", request.PartnerServiceID)
		}

		if request.CustomerNo != "123456789012345678" {
			t.Errorf("unexpected customerNo: %s", request.CustomerNo)
		}

		if request.VirtualAccountNo != "12345678123456789012345678" {
			t.Errorf("unexpected virtualAccountNo: %s", request.VirtualAccountNo)
		}

		if request.InquiryRequestID != "202010290000000000000000000001" {
			t.Errorf("unexpected inquiryRequestId: %s", request.InquiryRequestID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{
			"responseCode": "2002400",
			"responseMessage": "Successful",
			"virtualAccountData": {
				"paymentFlagStatus": "00",
				"paymentFlagReason": {
					"english": "Payment successful",
					"indonesia": "Pembayaran berhasil"
				},
				"partnerServiceId": "12345678",
				"customerNo": "123456789012345678",
				"virtualAccountNo": "12345678123456789012345678",
				"inquiryRequestId": "202010290000000000000000000001",
				"paymentRequestId": "202010290000000000000000000001",
				"paidAmount": {
					"value": "500000.00",
					"currency": "IDR"
				},
				"paidBills": "FF",
				"totalAmount": {
					"value": "500000.00",
					"currency": "IDR"
				},
				"trxDateTime": "2023-02-02T10:00:00+07:00",
				"transactionDate": "2023-02-02T10:00:00+07:00",
				"referenceNo": "12345678901",
				"paymentType": "1",
				"flagAdvise": "0",
				"billDetails": [
					{
						"billCode": "01",
						"billNo": "123456789012345678",
						"billName": "Electricity",
						"billShortName": "Electric",
						"billDescription": {
							"english": "Electricity Bill",
							"indonesia": "Tagihan Listrik"
						},
						"billSubCompany": "12345",
						"billAmount": {
							"value": "500000.00",
							"currency": "IDR"
						},
						"additionalInfo": {
							"customerType": "retail"
						},
						"billReferenceNo": "12345678901",
						"status": "00",
						"reason": {
							"english": "Payment successful",
							"indonesia": "Pembayaran berhasil"
						}
					}
				],
				"freeTexts": [
					{
						"english": "Thank you",
						"indonesia": "Terima kasih"
					}
				],
				"additionalInfo": {
					"customField": "custom-value"
				}
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

	service := NewVirtualAccountService(client)

	response, err := service.VirtualAccountInquiry(
		context.Background(),
		VirtualAccountInquiryRequest{
			PartnerServiceID: "12345678",
			CustomerNo:       "123456789012345678",
			VirtualAccountNo: "12345678123456789012345678",
			InquiryRequestID: "202010290000000000000000000001",
		},
		"28910000006578499987546738976812",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	if response.ResponseCode != "2002400" {
		t.Errorf("unexpected responseCode: %s", response.ResponseCode)
	}

	if response.ResponseMessage != "Successful" {
		t.Errorf("unexpected responseMessage: %s", response.ResponseMessage)
	}

	data := response.VirtualAccountData

	if data.PaymentFlagStatus != "00" {
		t.Errorf("unexpected paymentFlagStatus: %s", data.PaymentFlagStatus)
	}

	if data.PaymentFlagReason.English != "Payment successful" {
		t.Errorf(
			"unexpected paymentFlagReason english: %s",
			data.PaymentFlagReason.English,
		)
	}

	if data.PaymentFlagReason.Indonesia != "Pembayaran berhasil" {
		t.Errorf(
			"unexpected paymentFlagReason indonesia: %s",
			data.PaymentFlagReason.Indonesia,
		)
	}

	if data.PartnerServiceID != "12345678" {
		t.Errorf("unexpected partnerServiceId: %s", data.PartnerServiceID)
	}

	if data.CustomerNo != "123456789012345678" {
		t.Errorf("unexpected customerNo: %s", data.CustomerNo)
	}

	if data.VirtualAccountNo != "12345678123456789012345678" {
		t.Errorf("unexpected virtualAccountNo: %s", data.VirtualAccountNo)
	}

	if data.InquiryRequestID != "202010290000000000000000000001" {
		t.Errorf("unexpected inquiryRequestId: %s", data.InquiryRequestID)
	}

	if data.PaymentRequestID != "202010290000000000000000000001" {
		t.Errorf("unexpected paymentRequestId: %s", data.PaymentRequestID)
	}

	if data.PaidAmount.Value != "500000.00" {
		t.Errorf("unexpected paidAmount value: %s", data.PaidAmount.Value)
	}

	if data.PaidAmount.Currency != "IDR" {
		t.Errorf("unexpected paidAmount currency: %s", data.PaidAmount.Currency)
	}

	if data.PaidBills != "FF" {
		t.Errorf("unexpected paidBills: %s", data.PaidBills)
	}

	if data.TotalAmount.Value != "500000.00" {
		t.Errorf("unexpected totalAmount value: %s", data.TotalAmount.Value)
	}

	if data.TotalAmount.Currency != "IDR" {
		t.Errorf("unexpected totalAmount currency: %s", data.TotalAmount.Currency)
	}

	if data.TransactionDateTime != "2023-02-02T10:00:00+07:00" {
		t.Errorf("unexpected trxDateTime: %s", data.TransactionDateTime)
	}

	if data.TransactionDate != "2023-02-02T10:00:00+07:00" {
		t.Errorf("unexpected transactionDate: %s", data.TransactionDate)
	}

	if data.ReferenceNo != "12345678901" {
		t.Errorf("unexpected referenceNo: %s", data.ReferenceNo)
	}

	if data.PaymentType != "1" {
		t.Errorf("unexpected paymentType: %s", data.PaymentType)
	}

	if data.FlagAdvise != "0" {
		t.Errorf("unexpected flagAdvise: %s", data.FlagAdvise)
	}

	if len(data.BillDetails) != 1 {
		t.Fatalf("expected 1 bill detail, got %d", len(data.BillDetails))
	}

	bill := data.BillDetails[0]

	if bill.BillCode != "01" {
		t.Errorf("unexpected billCode: %s", bill.BillCode)
	}

	if bill.BillNo != "123456789012345678" {
		t.Errorf("unexpected billNo: %s", bill.BillNo)
	}

	if bill.BillName != "Electricity" {
		t.Errorf("unexpected billName: %s", bill.BillName)
	}

	if bill.BillShortName != "Electric" {
		t.Errorf("unexpected billShortName: %s", bill.BillShortName)
	}

	if bill.BillDescription.English != "Electricity Bill" {
		t.Errorf(
			"unexpected billDescription english: %s",
			bill.BillDescription.English,
		)
	}

	if bill.BillDescription.Indonesia != "Tagihan Listrik" {
		t.Errorf(
			"unexpected billDescription indonesia: %s",
			bill.BillDescription.Indonesia,
		)
	}

	if bill.BillSubCompany != "12345" {
		t.Errorf("unexpected billSubCompany: %s", bill.BillSubCompany)
	}

	if bill.BillAmount.Value != "500000.00" {
		t.Errorf("unexpected billAmount value: %s", bill.BillAmount.Value)
	}

	if bill.BillAmount.Currency != "IDR" {
		t.Errorf("unexpected billAmount currency: %s", bill.BillAmount.Currency)
	}

	if bill.BillReferenceNo != "12345678901" {
		t.Errorf("unexpected billReferenceNo: %s", bill.BillReferenceNo)
	}

	if bill.Status != "00" {
		t.Errorf("unexpected bill status: %s", bill.Status)
	}

	if bill.Reason.English != "Payment successful" {
		t.Errorf("unexpected bill reason english: %s", bill.Reason.English)
	}

	if bill.Reason.Indonesia != "Pembayaran berhasil" {
		t.Errorf("unexpected bill reason indonesia: %s", bill.Reason.Indonesia)
	}

	if len(data.FreeTexts) != 1 {
		t.Fatalf("expected 1 free text, got %d", len(data.FreeTexts))
	}

	if data.FreeTexts[0].English != "Thank you" {
		t.Errorf("unexpected free text english: %s", data.FreeTexts[0].English)
	}

	if data.FreeTexts[0].Indonesia != "Terima kasih" {
		t.Errorf("unexpected free text indonesia: %s", data.FreeTexts[0].Indonesia)
	}
}

func TestVirtualAccountService_VirtualAccountInquiry_Error(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		responseCode    string
		responseMessage string
	}{
		{
			name:            "bad request",
			statusCode:      http.StatusBadRequest,
			responseCode:    "4002600",
			responseMessage: "Bad Request",
		},
		{
			name:            "invalid field format",
			statusCode:      http.StatusBadRequest,
			responseCode:    "4002601",
			responseMessage: "Invalid Field Format accountNo",
		},
		{
			name:            "invalid mandatory field",
			statusCode:      http.StatusBadRequest,
			responseCode:    "4002602",
			responseMessage: "Invalid Mandatory Field accountNo",
		},
		{
			name:            "unauthorized",
			statusCode:      http.StatusUnauthorized,
			responseCode:    "4012600",
			responseMessage: "Unauthorized. [Reason]",
		},
		{
			name:            "invalid token",
			statusCode:      http.StatusUnauthorized,
			responseCode:    "4012601",
			responseMessage: "Invalid token (B2B)",
		},
		{
			name:            "transaction not found",
			statusCode:      http.StatusNotFound,
			responseCode:    "4042601",
			responseMessage: "Transaction Not Found",
		},
		{
			name:            "conflict",
			statusCode:      http.StatusConflict,
			responseCode:    "4092600",
			responseMessage: "Conflict",
		},
		{
			name:            "internal server error",
			statusCode:      http.StatusInternalServerError,
			responseCode:    "5002601",
			responseMessage: "Internal Server Error",
		},
		{
			name:            "timeout",
			statusCode:      http.StatusGatewayTimeout,
			responseCode:    "5042600",
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

			service := NewVirtualAccountService(client)

			response, err := service.VirtualAccountInquiry(
				context.Background(),
				VirtualAccountInquiryRequest{
					PartnerServiceID: "12345678",
					CustomerNo:       "123456789012345678",
					VirtualAccountNo: "12345678123456789012345678",
					InquiryRequestID: "202010290000000000000000000001",
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

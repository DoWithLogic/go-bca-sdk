package business_debit_card

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

func TestBusinessDebitCardService_CardStatusChange(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method %s, got %s", http.MethodPost, r.Method)
		}

		if r.URL.Path != "/openapi/bdc/v1.0/update-status" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if got := r.Header.Get("X-EXTERNAL-ID"); got != "external-123" {
			t.Errorf("expected X-EXTERNAL-ID %q, got %q", "external-123", got)
		}

		var request CardStatusChangeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		expectedRequest := CardStatusChangeRequest{
			CardNumber: "1234567890123456",
			BlockType:  BlockTypeFull,
			ReasonCode: ReasonCodeLost,
		}

		if request != expectedRequest {
			t.Errorf("unexpected request: %+v", request)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(CardStatusChangeResponse{
			ResponseCode:    "200US00",
			ResponseMessage: "Successful",
			ResponseData: CardStatusChangeResponseData{
				CardNumber: "1234567890123456",
				CardStatus: "B",
			},
		})
	}))
	defer server.Close()

	httpClient := server.Client()

	transportClient := transport.NewClient(
		httpClient,
		server.URL,
		mockAuthenticator{},
		transport.RetryConfig{
			MaxRetries: 0,
		},
	)

	service := &BusinessDebitCardService{
		transport: transportClient,
	}

	request := CardStatusChangeRequest{
		CardNumber: "1234567890123456",
		BlockType:  BlockTypeFull,
		ReasonCode: ReasonCodeLost,
	}

	response, err := service.CardStatusChange(
		context.Background(),
		request,
		"external-123",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.ResponseCode != "200US00" {
		t.Errorf(
			"expected response code %q, got %q",
			"200US00",
			response.ResponseCode,
		)
	}

	if response.ResponseMessage != "Successful" {
		t.Errorf(
			"expected response message %q, got %q",
			"Successful",
			response.ResponseMessage,
		)
	}

	if response.ResponseData.CardNumber != "1234567890123456" {
		t.Errorf(
			"expected card number %q, got %q",
			"1234567890123456",
			response.ResponseData.CardNumber,
		)
	}

	if response.ResponseData.CardStatus != "B" {
		t.Errorf(
			"expected card status %q, got %q",
			"B",
			response.ResponseData.CardStatus,
		)
	}
}

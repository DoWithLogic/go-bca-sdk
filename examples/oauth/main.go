package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	bca "github.com/DoWithLogic/go-bca-sdk"
	"github.com/DoWithLogic/go-bca-sdk/business_debit_card"
	bcaerrors "github.com/DoWithLogic/go-bca-sdk/errors"
)

func main() {
	client, err := bca.NewClient(
		bca.WithEnvironment(bca.Sandbox),
		bca.WithClientID("YOUR-CLIENT-ID"),
		bca.WithClientSecret("YOUR-CLIENT-SECRET"),
		bca.WithAPISecret("YOUR-API-SECRET"),
		bca.WithChannelID("YOUR-CHANNEL-ID"),
		bca.WithPartnerID("YOUR-PARTNER-ID"),
	)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.BusinessDebitCard.DebitCardInquiry(
		context.Background(),
		business_debit_card.DebitCardInquiryRequest{
			BankCardToken: "1234567890123456",
		},
		"01234506022024001",
	)
	if err != nil {
		var apiErr *bcaerrors.APIError
		if errors.As(err, &apiErr) {
			fmt.Println("HTTP Status:", apiErr.HTTPStatusCode)
			fmt.Println("Response Code:", apiErr.ResponseCode)
			fmt.Println("Response Message:", apiErr.ResponseMessage)
			log.Fatal(apiErr)
		}

		log.Fatal(err)
	}

	fmt.Println("Response Code:", response.ResponseCode)
	fmt.Println("Response Message:", response.ResponseMessage)
	fmt.Println("Response Data:", response.ResponseData)
}

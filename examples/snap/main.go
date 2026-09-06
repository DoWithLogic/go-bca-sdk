package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"

	bca "github.com/DoWithLogic/go-bca-sdk"
	"github.com/DoWithLogic/go-bca-sdk/account_information"
	bcaErrors "github.com/DoWithLogic/go-bca-sdk/errors"
)

func main() {
	yourPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	client, err := bca.NewClient(
		bca.WithEnvironment(bca.Sandbox),
		bca.WithClientID("YOUR-CLIENT-ID"),
		bca.WithClientSecret("YOUR-CLIENT-SECRET"),
		bca.WithChannelID("YOUR-CHANNEL-ID"),
		bca.WithPartnerID("YOUR-PARTNER-ID"),
		bca.WithSNAPAuth(yourPrivateKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.AccountInformation.BalanceInquiry(
		context.Background(),
		account_information.BalanceInquiryRequest{
			PartnerReferenceNo: "TEST-001",
			AccountNo:          "1234567890",
		},
		"external-id-001",
	)
	if err != nil {
		var apiErr *bcaErrors.APIError
		if errors.As(err, &apiErr) {
			fmt.Println("HTTP Status:", apiErr.HTTPStatusCode)
			fmt.Println("Response Code:", apiErr.ResponseCode)
			fmt.Println("Response Message:", apiErr.ResponseMessage)

			log.Fatal(apiErr.Error())

			return
		}

		log.Fatal(err)
	}

	fmt.Println("Response Code:", response.ResponseCode)
	fmt.Println("Response Message:", response.ResponseMessage)
	fmt.Println("Response Data:", response)
}

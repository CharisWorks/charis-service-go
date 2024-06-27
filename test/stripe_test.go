package test

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/charisworks/charisworks-service-go/strapi"
	"github.com/charisworks/charisworks-service-go/util"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/transferreversal"
)

func TestItemRegister(t *testing.T) {
	stripe.Key = "sk_test_51Nj1urA3bJzqElthGP4F3QjdR0SKk77E4pGHrsBAQEHia6lasXyujFOKXDyrodAxaE6PH6u2kNCVSdC5dBIRh82u00XqHQIZjM"

	print("Hello World")
	params := &stripe.PriceParams{
		Currency:    stripe.String(string(stripe.CurrencyJPY)),
		UnitAmount:  stripe.Int64(1000),
		ProductData: &stripe.PriceProductDataParams{Name: stripe.String("test item")},
	}
	result, err := price.New(params)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	t.Logf("Prices: %v", result)
}
func TestCreatePrice(t *testing.T) {
	stripe.Key = "sk_test_51Nj1urA3bJzqElthGP4F3QjdR0SKk77E4pGHrsBAQEHia6lasXyujFOKXDyrodAxaE6PH6u2kNCVSdC5dBIRh82u00XqHQIZjM"

	params := &stripe.PriceParams{
		Currency: stripe.String(string(stripe.CurrencyJPY)),
		ProductData: &stripe.PriceProductDataParams{
			ID:   stripe.String("test_2"),
			Name: stripe.String("test item"),
		},
		UnitAmount: stripe.Int64(1000),
	}
	result, err := price.New(params)
	if err != nil {

		t.Fatalf("Error: %v", err.Error())

	}

	t.Logf("Prices: %v", result)
	t.Logf("Prices: %v", result.ID)
}
func TestCreateCheckoutSession(t *testing.T) {
	now := time.Now()
	stripe.Key = "sk_test_51Nj1urA3bJzqElthGP4F3QjdR0SKk77E4pGHrsBAQEHia6lasXyujFOKXDyrodAxaE6PH6u2kNCVSdC5dBIRh82u00XqHQIZjM"
	params := &stripe.CheckoutSessionParams{
		ExpiresAt:  stripe.Int64(now.Add(30 * time.Minute).Unix()),
		SuccessURL: stripe.String("https://example.com/success"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("price_1PUOqgA3bJzqElth7EJRA5R2"),
				Quantity: stripe.Int64(2),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	result, err := session.New(params)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	t.Logf("Session: %v", result)
	t.Logf("Session: %v", result.ExpiresAt)
}

func TestRegisterPriceId(t *testing.T) {
	err := strapi.RegisterPriceId(1, "price_1PUOqgA3bJzqElth7EJRA5R2")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

}
func TestFetch(t *testing.T) {

	req, _ := http.NewRequest(string(GET), "https://example.com", nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+util.STRAPI_JWT)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf(err.Error())
	}
	log.Print(res)
}
func TestReverse(t *testing.T) {
	stripe.Key = "sk_test_51Nj1urA3bJzqElthGP4F3QjdR0SKk77E4pGHrsBAQEHia6lasXyujFOKXDyrodAxaE6PH6u2kNCVSdC5dBIRh82u00XqHQIZjM"

	reverseParams := &stripe.TransferReversalParams{
		ID: stripe.String("tr_1PWAyBA3bJzqElthWUSECLuS"),
	}
	transferResult, err := transferreversal.New(reverseParams)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(transferResult)
}

type httpMethod string

const (
	GET    httpMethod = "GET"
	POST   httpMethod = "POST"
	PUT    httpMethod = "PUT"
	DELETE httpMethod = "DELETE"
)

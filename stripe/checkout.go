package stripe

import (
	"os"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func CreateCheckout(priceId string, quantity int) (paymentUrl string, err error) {
	now := time.Now()
	stripe.Key = os.Getenv("STRIPE_SECRET_API_KEY")
	params := &stripe.CheckoutSessionParams{
		ExpiresAt:  stripe.Int64(now.Add(30 * time.Minute).Unix()),
		SuccessURL: stripe.String("https://example.com/success"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(int64(quantity)),
			},
		},
		BillingAddressCollection: stripe.String(string(stripe.CheckoutSessionBillingAddressCollectionRequired)),
		Mode:                     stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	result, err := session.New(params)
	if err != nil {
		panic(err)
	}
	return result.URL, nil
}

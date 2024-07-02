package stripe

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func CreateCheckout(priceId string, quantity int) (paymentUrl string, csId string, err error) {
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
		ShippingOptions: []*stripe.CheckoutSessionShippingOptionParams{
			{
				ShippingRateData: &stripe.CheckoutSessionShippingOptionShippingRateDataParams{
					DisplayName: stripe.String("送料"),
					Type:        stripe.String("fixed_amount"),
					FixedAmount: &stripe.CheckoutSessionShippingOptionShippingRateDataFixedAmountParams{
						Amount:   stripe.Int64(300),
						Currency: stripe.String(string(stripe.CurrencyJPY)),
					},
				},
			},
		},

		BillingAddressCollection: stripe.String(string(stripe.CheckoutSessionBillingAddressCollectionRequired)),
		Mode:                     stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	result, err := session.New(params)
	if err != nil {
		panic(err)
	}
	log.Print(result.ID)

	// 構造体をJSONにエンコード
	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// JSONデータを表示
	fmt.Println(string(jsonData))

	return result.URL, result.ID, nil
}

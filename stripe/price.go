package stripe

import (
	"github.com/charisworks/charisworks-service-go/util"
	"github.com/stripe/stripe-go/v76"
	_price "github.com/stripe/stripe-go/v76/price"
)

func CreatePrice(itemName string, price int) (priceId string, err error) {
	// Create a new price.
	stripe.Key = util.STRIPE_SECRET_API_KEY

	params := &stripe.PriceParams{
		Currency: stripe.String(string(stripe.CurrencyJPY)),
		ProductData: &stripe.PriceProductDataParams{
			Name: stripe.String(itemName),
		},
		UnitAmount: stripe.Int64(int64(price)),
	}
	result, err := _price.New(params)
	if err != nil {
		return "", err
	}
	return result.ID, nil

}
func ArchivePrice(priceId string) (err error) {
	// Deactivate a price.
	stripe.Key = util.STRIPE_SECRET_API_KEY

	_, err = _price.Update(priceId, &stripe.PriceParams{
		Active: stripe.Bool(false),
	})
	if err != nil {
		return err
	}
	return nil
}

package stripe

import (
	"fmt"

	"github.com/charisworks/charisworks-service-go/util"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/transfer"
	"github.com/stripe/stripe-go/v76/transferreversal"
)

func Transfer(amount float64, stripeAccountId string, transactionId string) (transferId string, err error) {
	stripe.Key = util.STRIPE_SECRET_API_KEY
	util.Logger(
		fmt.Sprintf(
			`
			Transfering... 
			amount: %v
			stripeID: %v
			transactionId: %v
			`, float64(amount)*(1-util.MARGIN-util.STRIPE_MARGIN), stripeAccountId, transactionId,
		),
	)

	params := &stripe.TransferParams{
		Amount:      stripe.Int64(int64(int(amount*(1-util.MARGIN)) + util.SHIPPING_FEE)),
		Currency:    stripe.String(string(stripe.CurrencyJPY)),
		Destination: stripe.String(stripeAccountId),
		Description: stripe.String(transactionId),
	}
	tr, err := transfer.New(params)
	if err != nil {
		return "", err
	}
	util.Logger(tr.ID)
	return tr.ID, nil
}

func ReverseTransfer(transferId string) (err error) {
	stripe.Key = util.STRIPE_SECRET_API_KEY
	util.Logger("Reversing transfer... \n ")
	reverseParams := &stripe.TransferReversalParams{
		ID: stripe.String(transferId),
	}
	transferResult, err := transferreversal.New(reverseParams)
	if err != nil {
		panic(err)
	}
	util.Logger(transferResult.ID)
	return
}

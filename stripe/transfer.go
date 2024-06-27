package stripe

import (
	"log"

	"github.com/charisworks/charisworks-service-go/util"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/transfer"
	"github.com/stripe/stripe-go/v76/transferreversal"
)

func Transfer(amount float64, stripeAccountId string, transactionId string) (transferId string, err error) {
	stripe.Key = util.STRIPE_SECRET_API_KEY
	log.Print("Transfering... \n amount: ", float64(amount)*(1-util.MARGIN), "\n stripeID: ", stripeAccountId, "\n transactionId: ", transactionId)

	params := &stripe.TransferParams{
		Amount:      stripe.Int64(int64(amount * (1 - util.MARGIN))),
		Currency:    stripe.String(string(stripe.CurrencyJPY)),
		Destination: stripe.String(stripeAccountId),
		Description: stripe.String(transactionId),
	}
	tr, err := transfer.New(params)
	if err != nil {
		return "", err
	}
	log.Print(tr.ID)
	return tr.ID, nil
}

func ReverseTransfer(transferId string) (err error) {
	stripe.Key = util.STRIPE_SECRET_API_KEY
	log.Print("Reversing transfer... \n ")
	reverseParams := &stripe.TransferReversalParams{
		ID: stripe.String(transferId),
	}
	transferResult, err := transferreversal.New(reverseParams)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(transferResult)
	return
}

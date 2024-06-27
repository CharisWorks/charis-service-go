package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/charisworks/charisworks-service-go/strapi"
	_stripe "github.com/charisworks/charisworks-service-go/stripe"
	"github.com/charisworks/charisworks-service-go/util"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

func (h *Handler) SetupStripeEventHandler() {
	h.Router.POST("/webhooks/stripe", stripeWebhookMiddleware(), h.StripeEventHandler)
}

func stripeWebhookMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const MaxBodyBytes = int64(65536)
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)

		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
			ctx.Abort()
			return
		}
		// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key
		// You can find your endpoint's secret in your webhook settings
		endpointSecret := util.STRIPE_API_KEY
		event, err := webhook.ConstructEvent(body, ctx.Request.Header.Get("Stripe-Signature"), endpointSecret)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
			ctx.Abort() // Return a 400 error on a bad signature
			return
		}
		ctx.Set("Event", event)
		ctx.Next()
	}
}

func (h *Handler) StripeEventHandler(ctx *gin.Context) {
	event := ctx.MustGet("Event").(stripe.Event)
	// 構造体をJSONにエンコード
	jsonData, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// JSONデータを表示
	fmt.Println(string(jsonData))

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		billing := sanitizeNill(event.Data.Object["customer_details"].(map[string]interface{}))
		address := sanitizeNill(event.Data.Object["customer_details"].(map[string]interface{})["address"].(map[string]interface{}))

		// 構造体をJSONにエンコード
		transaction, err := strapi.GetTransactionById(event.Data.Object["id"].(string))
		fmt.Printf(`
*************************************************
CheckoutSession was completed!
transactionId: %s
****Customer Infomation****
state: %s
city: %s
line1: %s
line2: %s
postal_code: %s
email: %s
name: %s
phone: %s
****Transaction Information****
ItemId: %d
Item Name: %s
Quantity: %d
*************************************************
		`,
			strconv.Itoa(transaction.Data[0].ID),
			address["state"],
			address["city"],
			address["line1"],
			address["line2"],
			address["postal_code"],
			billing["email"],
			billing["name"],
			billing["phone"],
			transaction.Data[0].Attributes.Item.Data.Id,
			transaction.Data[0].Attributes.Item.Data.Attributes.Name,
			transaction.Data[0].Attributes.Quantity,
		)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := strapi.CheckoutSessionDetailRegister(
			strconv.Itoa(transaction.Data[0].ID),
			address["state"],
			address["city"],
			address["line1"],
			address["line2"],
			address["postal_code"],
			billing["email"],
			billing["name"],
			billing["phone"],
			event.Data.Object["payment_intent"].(string),
		); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := strapi.ReducePreStock(transaction.Data[0].Attributes.Item.Data.Id, transaction.Data[0].Attributes.Quantity); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		item, err := strapi.GetItem(transaction.Data[0].Attributes.Item.Data.Id)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		trId, err := _stripe.Transfer(event.Data.Object["amount_total"].(float64), item.Data.Attributes.Worker.Data.Attributes.StripeAccountID, transaction.Data[0].Attributes.TransactionID)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := strapi.CheckoutSessionTransferRegister(strconv.Itoa(transaction.Data[0].ID), trId); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
	case stripe.EventTypeCheckoutSessionExpired:
		fmt.Println("CheckoutSession was expired!")
		transaction, err := strapi.GetTransactionById(event.Data.Object["id"].(string))
		log.Print("got transaction: ", transaction)

		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := strapi.ReturnPreStock(transaction.Data[0].Attributes.Item.Data.Id, transaction.Data[0].Attributes.Quantity); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err = strapi.CheckoutSessionStatusRegister(strconv.Itoa(transaction.Data[0].ID), strapi.Cancelled); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	case stripe.EventTypeChargeRefunded:
		fmt.Println("Charge was refunded!")
		transaction, err := strapi.GetTransactionByPaymentIntent(event.Data.Object["payment_intent"].(string))
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = _stripe.ReverseTransfer(transaction.Data[0].Attributes.TransferID.(string))
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		log.Print("返金完了！！！！")
	}
	ctx.JSON(http.StatusOK, gin.H{"received": true})
}

func RegisterTransaction(event stripe.Event) {
	// Register a transfer.

}
func sanitizeNill(m map[string]interface{}) map[string]string {
	n := make(map[string]string)
	for c := range m {
		if name, ok := m[c].(string); ok {
			n[c] = name
		} else {
			n[c] = ""
		}
	}

	return n
}

package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"

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
		endpointSecret := os.Getenv("STRIPE_API_KEY")
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
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		fmt.Println("CheckoutSession was completed!")
	case stripe.EventTypeCheckoutSessionExpired:
		fmt.Println("CheckoutSession was expired!")
	case stripe.EventTypeChargeRefunded:
		fmt.Println("Charge was refunded!")
	}

	ctx.JSON(http.StatusOK, gin.H{"received": true})
}

package handler

import (
	"log"
	"net/http"

	"github.com/charisworks/charisworks-service-go/strapi"
	"github.com/charisworks/charisworks-service-go/stripe"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SetupCheckoutEventHandler() {
	h.Router.POST("/checkout", h.CheckoutHandler)
}

func (h *Handler) CheckoutHandler(ctx *gin.Context) {
	payload := transactionPayload{}
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	url, err := RegisterPendingTransaction(payload.ItemId, payload.Quantity)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"url": url})
}

func RegisterPendingTransaction(itemId int, quantity int) (redirectUrl string, err error) {
	// Register a pending transaction.
	if err := strapi.ShiftStock(itemId, quantity); err != nil {
		return "", err
	}
	priceId, err := strapi.GetItem(itemId)
	if err != nil {
		return "", err
	}
	url, err := stripe.CreateCheckout(priceId.Data.Attributes.PriceId, quantity)
	if err != nil {
		return "", err
	}
	return url, nil
}

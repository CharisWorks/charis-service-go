package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/charisworks/charisworks-service-go/meilisearch"
	"github.com/charisworks/charisworks-service-go/strapi"
	_stripe "github.com/charisworks/charisworks-service-go/stripe"
	"github.com/charisworks/charisworks-service-go/util"
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
	url, err := registerPendingTransaction(payload.ItemId, payload.Quantity)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": url})
}

func registerPendingTransaction(itemId int, quantity int) (redirectUrl string, err error) {
	// Register a pending transaction.
	priceId, err := strapi.GetItem(itemId)
	if err != nil {
		return "", err
	}
	if priceId.Data.Attributes.PublishedAt == "" {
		err = util.NewError("this item is not published")
		return "", err
	}
	if err := strapi.ShiftStock(itemId, quantity); err != nil {
		return "", err
	}
	url, csId, err := _stripe.CreateCheckout(priceId.Data.Attributes.PriceId, quantity)
	if err != nil {
		return "", err
	}
	if err := strapi.TransactionRegister(csId, itemId, quantity, strapi.Pending); err != nil {
		return "", err
	}
	item, err := meilisearch.GetItemByID(strconv.Itoa(itemId))
	if err != nil {
		return "", err
	}
	item.Stock -= quantity
	items := []meilisearch.Item{}
	items = append(items, *item)
	if err := meilisearch.RegisterItemToMeilisearch(items); err != nil {
		return "", err
	}
	return url, nil
}

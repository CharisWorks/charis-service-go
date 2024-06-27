package handler

import (
	"log"
	"net/http"

	"github.com/charisworks/charisworks-service-go/strapi"
	"github.com/charisworks/charisworks-service-go/stripe"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SetupStrapiEventHandler() {
	h.Router.POST("/webhooks/strapi", strapiWebhookMiddleware(), h.StrapiEventHandler)
	h.Router.POST("/webhooks/strapi/create", strapiWebhookMiddleware(), h.StrapiCreateEventHandler)
	h.Router.POST("/webhooks/strapi/itemdelete", strapiWebhookMiddleware(), h.StrapiItemDeleteHandler)
}

func strapiWebhookMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func (h *Handler) StrapiEventHandler(ctx *gin.Context) {
	event := &strapi.Event{}
	err := ctx.BindJSON(&event)
	if err != nil {
		log.Print(event)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Print(event)
	// body, err := io.ReadAll(ctx.Request.Body)
	// if err != nil {
	// 	ctx.AbortWithStatus(http.StatusBadRequest)
	// 	return
	// }

	// // リクエストメソッドとURLを表示
	// fmt.Printf("Received %s request to %s\n", ctx.Request.Method, ctx.Request.RequestURI)

	// // リクエストボディを表示
	// fmt.Printf("Request body: %s\n", body)

	// 200 OKを返す
	ctx.JSON(http.StatusOK, gin.H{"received": true})
}
func (h *Handler) StrapiCreateEventHandler(ctx *gin.Context) {
	event := &strapi.ItemEvent{}
	err := ctx.BindJSON(&event)
	if err != nil {
		log.Print(event)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Print(event)
	if event.Model == strapi.ItemModel {
		priceId, err := stripe.CreatePrice(event.Entry.Name, event.Entry.Price)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.Print("successfully registered item: ", priceId)
		err = strapi.RegisterPriceId(event.Entry.ID, priceId)
		if err != nil {
			log.Print(err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	// 200 OKを返す
	ctx.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *Handler) StrapiItemDeleteHandler(ctx *gin.Context) {
	event := &strapi.Event{}
	err := ctx.BindJSON(&event)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if event.EventName != strapi.Delete {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	item := strapi.ItemEvent{}
	err = stripe.ArchivePrice(item.Entry.PriceId)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Print("successfully archived item: ", item.Entry.ID)

	// 200 OKを返す
	ctx.JSON(http.StatusOK, gin.H{"received": true})
}

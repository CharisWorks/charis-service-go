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
	h.Router.POST("/webhooks/strapi/itemregister", strapiWebhookMiddleware(), h.StrapiItemRegisterHandler)
	h.Router.POST("/webhooks/strapi/itemdelete", strapiWebhookMiddleware(), h.StrapiItemDeleteHandler)

}

func strapiWebhookMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func (h *Handler) StrapiEventHandler(ctx *gin.Context) {
	event := &strapi.Item{}
	err := ctx.BindJSON(&event)
	if err != nil {
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
func (h *Handler) StrapiItemRegisterHandler(ctx *gin.Context) {
	item := &strapi.Item{}
	err := ctx.BindJSON(&item)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	priceId, err := stripe.CreatePrice(item.Entry.Name, item.Entry.Price)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Print(priceId)
	err = strapi.RegisterPriceId(item.Entry.ID, priceId)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// 200 OKを返す
	ctx.JSON(http.StatusOK, gin.H{"received": true})
}

func (h *Handler) StrapiItemDeleteHandler(ctx *gin.Context) {
	item := &strapi.Item{}
	err := ctx.BindJSON(&item)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = stripe.ArchivePrice(item.Entry.PriceId)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 200 OKを返す
	ctx.JSON(http.StatusOK, gin.H{"received": true})
}

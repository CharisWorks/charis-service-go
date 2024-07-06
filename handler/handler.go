package handler

import (
	"github.com/charisworks/charisworks-service-go/images"

	"github.com/gin-gonic/gin"
)

var r2conns = images.R2Conns{}

type Handler struct {
	Router *gin.Engine
}

func NewHandler(router *gin.Engine) *Handler {
	r2conns.Init()
	return &Handler{
		Router: router,
	}
}
func (h *Handler) SetupStripeEventHandler() {
	h.Router.POST("/webhooks/stripe", stripeWebhookMiddleware(), h.stripeEventHandler)
}
func (h *Handler) SetupStrapiEventHandler() {
	h.Router.POST("/webhooks/strapi", strapiWebhookMiddleware(), h.StrapiEventHandler)
}

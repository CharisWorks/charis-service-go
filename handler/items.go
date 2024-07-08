package handler

import (
	"log"
	"net/http"

	"github.com/charisworks/charisworks-service-go/meilisearch"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SetupItmeHandler() {
	h.Router.GET("/item/:item_id", h.ItemHandler)
}

func (h *Handler) ItemHandler(ctx *gin.Context) {
	itemId := ctx.Param("item_id")
	item, err := meilisearch.GetItemByID(itemId)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, &item)
}

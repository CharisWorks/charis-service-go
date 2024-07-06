package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/charisworks/charisworks-service-go/meilisearch"
	"github.com/charisworks/charisworks-service-go/strapi"
	"github.com/charisworks/charisworks-service-go/stripe"
	"github.com/charisworks/charisworks-service-go/util"

	"github.com/gin-gonic/gin"
)

func strapiWebhookMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func (h *Handler) StrapiEventHandler(ctx *gin.Context) {

	event := &strapi.Event{}
	if err := ctx.ShouldBindBodyWithJSON(&event); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	util.Logger(
		fmt.Sprintf(
			`
			*************************************************
			Strapi Event was received!
			Model: %v
			EventName: %v
			*************************************************
			`, event.Model, event.EventName,
		),
	)

	if event.Model == strapi.ItemModel {
		switch event.EventName {
		case strapi.Update:
			itemEvent := &strapi.ItemEvent{}
			if err := ctx.ShouldBindBodyWithJSON(&itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			if err := itemUpdateHandler(itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		case strapi.Publish:
			itemEvent := &strapi.ItemEvent{}
			if err := ctx.ShouldBindBodyWithJSON(&itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			if err := itemPublishHandler(itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		case strapi.Unpublish:
			itemEvent := &strapi.ItemEvent{}
			if err := ctx.ShouldBindBodyWithJSON(&itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			if err := itemUnpublishHandler(itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		case strapi.Delete:
			itemEvent := &strapi.ItemEvent{}
			if err := ctx.ShouldBindBodyWithJSON(&itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			if err := itemDeleteHandler(itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}
	}

	if event.Model == strapi.TransactionModel {
		switch event.EventName {
		case strapi.Update:
			transactionEvent := &strapi.TransactionEvent{}
			if err := ctx.ShouldBindBodyWithJSON(&transactionEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			if transactionEvent.Entry.Status == strapi.ShippedTransaction {
				//transactionEventの中身を取得
				transaction, err := strapi.GetTransactionById(transactionEvent.Entry.TransactionID)
				if err != nil {
					ctx.AbortWithStatus(http.StatusBadRequest)
					return
				}
				err = shippingHandler(transaction)
				if err != nil {
					ctx.AbortWithStatus(http.StatusBadRequest)
					return
				}
			}
		}
	}

	// 200 OKを返す
	ctx.JSON(http.StatusOK, gin.H{"received": true})
}

func itemPublishHandler(itemEvent *strapi.ItemEvent) (err error) {
	item, err := strapi.GetItem(itemEvent.Entry.ID)
	if err != nil {
		return err
	}
	if item.Data.Attributes.PublishedAt == "" {
		return nil
	}
	if item.Data.Attributes.PriceId == "" {
		priceId, err := stripe.CreatePrice(itemEvent.Entry.Name, itemEvent.Entry.Price)
		if err != nil {
			return err
		}

		err = strapi.RegisterPriceId(itemEvent.Entry.ID, priceId)
		if err != nil {
			return err
		}
	}
	Image := new([]meilisearch.Images)
	for _, img := range item.Data.Attributes.Images.Data {

		i := meilisearch.Images{}
		if img.Attributes.Formats.Small.Url != "" {
			i.SmallUrl = util.IMAGES_URL + "/" + strconv.Itoa(item.Data.Id) + "/" + img.Attributes.Formats.Small.Hash
		}
		if img.Attributes.Formats.Medium.Url != "" {
			i.MediumUrl = util.IMAGES_URL + "/" + strconv.Itoa(item.Data.Id) + "/" + img.Attributes.Formats.Medium.Hash
		}
		if img.Attributes.Formats.Large.Url != "" {
			i.LargeUrl = util.IMAGES_URL + "/" + strconv.Itoa(item.Data.Id) + "/" + img.Attributes.Formats.Large.Hash
		}
		*Image = append(*Image, i)
	}
	if err := meilisearch.RegisterItemToMeilisearch(
		[]meilisearch.Item{
			{
				ID:                strconv.Itoa(item.Data.Id),
				ItemName:          item.Data.Attributes.Name,
				Price:             item.Data.Attributes.Price,
				Stock:             item.Data.Attributes.Stock,
				Description:       item.Data.Attributes.Description,
				Genre:             item.Data.Attributes.Genre,
				CreatedAt:         item.Data.Attributes.CreatedAt,
				UpdatedAt:         item.Data.Attributes.UpdatedAt,
				PublishedAt:       item.Data.Attributes.PublishedAt,
				Worker:            item.Data.Attributes.Worker.Data.Attributes.UserName,
				WorkerDescription: item.Data.Attributes.Worker.Data.Attributes.Description,
				ThumbnailUrl:      util.IMAGES_URL + "/" + strconv.Itoa(item.Data.Id) + "/" + item.Data.Attributes.Images.Data[0].Attributes.Formats.Thumbnail.Hash,
				Images:            *Image,
				Tags:              []string{item.Data.Attributes.Tag.Color, item.Data.Attributes.Tag.Size},
			},
		},
	); err != nil {
		return err
	}
	return nil
}
func itemUpdateHandler(itemEvent *strapi.ItemEvent) (err error) {
	item, err := strapi.GetItem(itemEvent.Entry.ID)
	if err != nil {
		return err
	}
	if item.Data.Attributes.PublishedAt == "" {
		return nil
	}
	util.Logger(
		fmt.Sprintf(`
	*************************************************
	Item was updated!
	ID: %v
	Name: %v
	Image: %v
	*************************************************
	`,
			item.Data.Id, item.Data.Attributes.Name, item.Data.Attributes.Images.Data[0].Attributes.Formats.Medium.Url),
	)
	if err := r2conns.UploadImage("."+item.Data.Attributes.Images.Data[0].Attributes.Formats.Thumbnail.Url, strconv.Itoa(item.Data.Id)+"/"+item.Data.Attributes.Images.Data[0].Attributes.Formats.Thumbnail.Hash); err != nil {
		return err
	}
	for _, img := range item.Data.Attributes.Images.Data {
		if img.Attributes.Formats.Small.Url != "" {
			if err := r2conns.UploadImage("."+img.Attributes.Formats.Small.Url, strconv.Itoa(item.Data.Id)+"/"+img.Attributes.Formats.Small.Hash); err != nil {
				return err
			}
		}
		if img.Attributes.Formats.Medium.Url != "" {
			if err := r2conns.UploadImage("."+img.Attributes.Formats.Medium.Url, strconv.Itoa(item.Data.Id)+"/"+img.Attributes.Formats.Medium.Hash); err != nil {
				return err
			}
		}
		if img.Attributes.Formats.Large.Url != "" {
			if err := r2conns.UploadImage("."+img.Attributes.Formats.Large.Url, strconv.Itoa(item.Data.Id)+"/"+img.Attributes.Formats.Large.Hash); err != nil {
				return err
			}
		}
	}
	util.Logger(
		`
*************************************************
Images were uploaded!
*************************************************
			`)
	Image := new([]meilisearch.Images)
	for _, img := range item.Data.Attributes.Images.Data {

		i := meilisearch.Images{}
		if img.Attributes.Formats.Small.Url != "" {
			i.SmallUrl = util.IMAGES_URL + "/" + strconv.Itoa(item.Data.Id) + "/" + img.Attributes.Formats.Small.Hash
		}
		if img.Attributes.Formats.Medium.Url != "" {
			i.MediumUrl = util.IMAGES_URL + "/" + strconv.Itoa(item.Data.Id) + "/" + img.Attributes.Formats.Medium.Hash
		}
		if img.Attributes.Formats.Large.Url != "" {
			i.LargeUrl = util.IMAGES_URL + "/" + strconv.Itoa(item.Data.Id) + "/" + img.Attributes.Formats.Large.Hash
		}
		*Image = append(*Image, i)
	}

	if err := meilisearch.RegisterItemToMeilisearch(
		[]meilisearch.Item{
			{
				ID:                strconv.Itoa(item.Data.Id),
				ItemName:          item.Data.Attributes.Name,
				Price:             item.Data.Attributes.Price,
				Stock:             item.Data.Attributes.Stock,
				Description:       item.Data.Attributes.Description,
				Genre:             item.Data.Attributes.Genre,
				CreatedAt:         item.Data.Attributes.CreatedAt,
				UpdatedAt:         item.Data.Attributes.UpdatedAt,
				PublishedAt:       item.Data.Attributes.PublishedAt,
				Worker:            item.Data.Attributes.Worker.Data.Attributes.UserName,
				WorkerDescription: item.Data.Attributes.Worker.Data.Attributes.Description,
				ThumbnailUrl:      util.IMAGES_URL + "/" + strconv.Itoa(item.Data.Id) + "/" + item.Data.Attributes.Images.Data[0].Attributes.Formats.Thumbnail.Hash,
				Images:            *Image,
				Tags:              []string{item.Data.Attributes.Tag.Color, item.Data.Attributes.Tag.Size},
			},
		},
	); err != nil {
		return err
	}
	return nil
}

func itemUnpublishHandler(itemEvent *strapi.ItemEvent) (err error) {
	if err := meilisearch.DeleteItemFromMeilisearch(strconv.Itoa(itemEvent.Entry.ID)); err != nil {
		return err
	}
	images, err := r2conns.GetImages(strconv.Itoa(itemEvent.Entry.ID))
	if err != nil {
		return err
	}
	for _, img := range images {
		if err := r2conns.DeleteImage(img); err != nil {
			return err
		}
	}
	return nil
}

func itemDeleteHandler(itemEvent *strapi.ItemEvent) (err error) {
	if err := stripe.ArchivePrice(itemEvent.Entry.PriceId); err != nil {
		return err
	}
	if err := meilisearch.DeleteItemFromMeilisearch(strconv.Itoa(itemEvent.Entry.ID)); err != nil {
		return err
	}
	images, err := r2conns.GetImages(strconv.Itoa(itemEvent.Entry.ID))
	if err != nil {
		return err
	}
	for _, img := range images {
		if err := r2conns.DeleteImage(img); err != nil {
			return err
		}
	}
	return nil
}

package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/charisworks/charisworks-service-go/meilisearch"
	"github.com/charisworks/charisworks-service-go/strapi"
	"github.com/charisworks/charisworks-service-go/stripe"
	"github.com/charisworks/charisworks-service-go/util"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SetupStrapiEventHandler() {
	h.Router.POST("/webhooks/strapi", strapiWebhookMiddleware(), h.StrapiEventHandler)
}

func strapiWebhookMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

// func (h *Handler) StrapiEventHandler(ctx *gin.Context) {
// 	event := &strapi.Event{}
// 	err := ctx.BindJSON(&event)
// 	if err != nil {
// 		log.Print(event)
// 		ctx.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
// 	log.Print(event)
// 	// body, err := io.ReadAll(ctx.Request.Body)
// 	// if err != nil {
// 	// 	ctx.AbortWithStatus(http.StatusBadRequest)
// 	// 	return
// 	// }

// 	// // リクエストメソッドとURLを表示
// 	// fmt.Printf("Received %s request to %s\n", ctx.Request.Method, ctx.Request.RequestURI)

// 	// // リクエストボディを表示
// 	// fmt.Printf("Request body: %s\n", body)

//		// 200 OKを返す
//		ctx.JSON(http.StatusOK, gin.H{"received": true})
//	}
func (h *Handler) StrapiEventHandler(ctx *gin.Context) {

	event := &strapi.Event{}
	if err := ctx.ShouldBindBodyWithJSON(&event); err != nil {
		log.Print(event)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Printf(`
	*************************************************
	Strapi Event was received!
	Model: %v
	EventName: %v
	*************************************************
	`, event.Model, event.EventName)

	if event.Model == strapi.ItemModel {
		switch event.EventName {
		case strapi.Create:
			// itemEvent := &strapi.ItemEvent{}
			// if err := ctx.ShouldBindBodyWithJSON(&itemEvent); err != nil {
			// 	ctx.AbortWithStatus(http.StatusBadRequest)
			// }
			// if err := ItemCreateHandler(itemEvent); err != nil {
			// 	ctx.AbortWithStatus(http.StatusBadRequest)
			// 	return
			// }
		case strapi.Update:
			itemEvent := &strapi.ItemEvent{}
			if err := ctx.ShouldBindBodyWithJSON(&itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			if err := ItemUpdateHandler(itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		case strapi.Publish:
			itemEvent := &strapi.ItemEvent{}
			if err := ctx.ShouldBindBodyWithJSON(&itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			if err := ItemPublishHandler(itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		case strapi.Unpublish:
			itemEvent := &strapi.ItemEvent{}
			if err := ctx.ShouldBindBodyWithJSON(&itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			if err := ItemUnpublishHandler(itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		case strapi.Delete:
			itemEvent := &strapi.ItemEvent{}
			if err := ctx.ShouldBindBodyWithJSON(&itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			if err := ItemDeleteHandler(itemEvent); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}
	}

	if event.Model == strapi.TransactionModel {
		switch event.EventName {
		case strapi.Update:
			log.Print("transactionEvent")
			transactionEvent := &strapi.TransactionEvent{}
			if err := ctx.ShouldBindBodyWithJSON(&transactionEvent); err != nil {
				log.Print(err)
				ctx.AbortWithStatus(http.StatusBadRequest)
			}
			log.Print("transactionEvent: ", transactionEvent)
			if transactionEvent.Entry.Status == strapi.ShippedTransaction {
				//transactionEventの中身を取得
				transaction, err := strapi.GetTransactionById(transactionEvent.Entry.TransactionID)
				if err != nil {
					ctx.AbortWithStatus(http.StatusBadRequest)
					return
				}
				err = ShippingHandler(transaction)
				if err != nil {
					ctx.AbortWithStatus(http.StatusBadRequest)
					return
				}
				log.Print("successfully registered transaction: ", transaction.Data[0].ID)
			}
		}
	}

	// 200 OKを返す
	ctx.JSON(http.StatusOK, gin.H{"received": true})
}

// func ItemCreateHandler(itemEvent *strapi.ItemEvent) (err error) {
// 	priceId, err := stripe.CreatePrice(itemEvent.Entry.Name, itemEvent.Entry.Price)
// 	if err != nil {
// 		return err
// 	}
// 	log.Print("successfully registered item: ", priceId)
// 	err = strapi.RegisterPriceId(itemEvent.Entry.ID, priceId)
// 	if err != nil {
// 		return err
// 	}
// 	item, err := strapi.GetItem(itemEvent.Entry.ID)
// 	if err != nil {
// 		return err
// 	}
// 	if item.Data.Attributes.PublishedAt == "" {
// 		return nil
// 	}
// 	if err := meilisearch.RegisterItemToMeilisearch(
// 		[]meilisearch.Item{
// 			{
// 				ID:          strconv.Itoa(item.Data.Id),
// 				ItemName:    item.Data.Attributes.Name,
// 				Price:       item.Data.Attributes.Price,
// 				Stock:       item.Data.Attributes.Stock,
// 				Description: item.Data.Attributes.Description,
// 				Genre:       item.Data.Attributes.Genre,
// 				CreatedAt:   item.Data.Attributes.CreatedAt,
// 				UpdatedAt:   item.Data.Attributes.UpdatedAt,
// 				PublishedAt: item.Data.Attributes.PublishedAt,
// 				Worker:      item.Data.Attributes.Worker.Data.Attributes.UserName,
// 			},
// 		},
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }

func ItemPublishHandler(itemEvent *strapi.ItemEvent) (err error) {
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

		log.Print("successfully registered item: ", priceId)
		err = strapi.RegisterPriceId(itemEvent.Entry.ID, priceId)
		if err != nil {
			return err
		}
	}

	if err := meilisearch.RegisterItemToMeilisearch(
		[]meilisearch.Item{
			{
				ID:          strconv.Itoa(item.Data.Id),
				ItemName:    item.Data.Attributes.Name,
				Price:       item.Data.Attributes.Price,
				Stock:       item.Data.Attributes.Stock,
				Description: item.Data.Attributes.Description,
				Genre:       item.Data.Attributes.Genre,
				CreatedAt:   item.Data.Attributes.CreatedAt,
				UpdatedAt:   item.Data.Attributes.UpdatedAt,
				PublishedAt: item.Data.Attributes.PublishedAt,
				Worker:      item.Data.Attributes.Worker.Data.Attributes.UserName,
			},
		},
	); err != nil {
		return err
	}
	return nil
}
func ItemUpdateHandler(itemEvent *strapi.ItemEvent) (err error) {
	item, err := strapi.GetItem(itemEvent.Entry.ID)
	if err != nil {
		return err
	}
	if item.Data.Attributes.PublishedAt == "" {
		return nil
	}
	log.Printf(`
	*************************************************
	Item was updated!
	ID: %v
	Name: %v
	Image: %v
	*************************************************
	`,
		item.Data.Id, item.Data.Attributes.Name, item.Data.Attributes.Images.Data[0].Attributes.Formats.Medium.Url)
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
	log.Println(`
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
				ID:           strconv.Itoa(item.Data.Id),
				ItemName:     item.Data.Attributes.Name,
				Price:        item.Data.Attributes.Price,
				Stock:        item.Data.Attributes.Stock,
				Description:  item.Data.Attributes.Description,
				Genre:        item.Data.Attributes.Genre,
				CreatedAt:    item.Data.Attributes.CreatedAt,
				UpdatedAt:    item.Data.Attributes.UpdatedAt,
				PublishedAt:  item.Data.Attributes.PublishedAt,
				Worker:       item.Data.Attributes.Worker.Data.Attributes.UserName,
				ThumbnailUrl: util.IMAGES_URL + "/" + strconv.Itoa(item.Data.Id) + "/" + item.Data.Attributes.Images.Data[0].Attributes.Formats.Thumbnail.Hash,
				Images:       *Image,
				Tags:         []string{item.Data.Attributes.Tag.Color, item.Data.Attributes.Tag.Size},
			},
		},
	); err != nil {
		return err
	}
	return nil
}

func ItemUnpublishHandler(itemEvent *strapi.ItemEvent) (err error) {
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

func ItemDeleteHandler(itemEvent *strapi.ItemEvent) (err error) {
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

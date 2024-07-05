package strapi

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/charisworks/charisworks-service-go/util"
)

func RegisterPriceId(itemId int, priceId string) error {
	// Register the price id.
	putRest, _ := json.Marshal(map[string]interface{}{
		"data": map[string]string{
			"price_id": priceId,
		},
	})
	res, err := requestToStrapi(PUT, "/items/"+strconv.Itoa(itemId), putRest)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return nil
}
func RegisterStockAndPrestock(itemId int, prestock int, stock int) error {
	// Register the price id.
	putRest, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"pre_stock": prestock,
			"stock":     stock,
		},
	})
	res, err := requestToStrapi(PUT, "/items/"+strconv.Itoa(itemId), putRest)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func GetItem(itemId int) (*Item, error) {
	// Get the price id.
	res, err := requestToStrapi(GET, "/items/"+strconv.Itoa(itemId)+"?populate[0]=worker&populate[1]=images&populate[2]=tag", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var item Item
	if err := json.NewDecoder(res.Body).Decode(&item); err != nil {
		return nil, err
	}

	return &item, nil
}

func ShiftStock(itemId int, quantity int) error {
	item, err := GetItem(itemId)
	if err != nil {
		return err
	}
	log.Print(item)
	if item.Data.Attributes.Stock < quantity {
		err = util.NewError("Not enough stock")
		return err
	}
	err = RegisterStockAndPrestock(itemId, item.Data.Attributes.PreStock+quantity, item.Data.Attributes.Stock-quantity)
	if err != nil {
		return err
	}
	return nil
}

func ReducePreStock(itemId int, quantity int) error {
	item, err := GetItem(itemId)
	if err != nil {
		return err
	}
	log.Printf(`
****************
itemid: %d
prestock: %d
quantity: %d
****************
	`, itemId, item.Data.Attributes.PreStock, quantity)
	err = RegisterStockAndPrestock(itemId, item.Data.Attributes.PreStock-quantity, item.Data.Attributes.Stock)
	if err != nil {
		return err
	}
	return nil
}
func ReturnPreStock(itemId int, quantity int) error {
	item, err := GetItem(itemId)
	if err != nil {
		return err
	}

	err = RegisterStockAndPrestock(itemId, item.Data.Attributes.PreStock-quantity, item.Data.Attributes.Stock+quantity)
	if err != nil {
		return err
	}
	return nil
}

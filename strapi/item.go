package strapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
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

	req, _ := http.NewRequest("PUT", "http://strapi:1337/api/items/"+strconv.Itoa(itemId), bytes.NewBuffer(putRest))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+os.Getenv("STRAPI_JWT"))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
func RegisterPrestock(itemId int, prestock int, stock int) error {
	// Register the price id.
	putRest, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"pre_stock": prestock,
			"stock":     stock,
		},
	})
	req, _ := http.NewRequest("PUT", "http://strapi:1337/api/items/"+strconv.Itoa(itemId), bytes.NewBuffer(putRest))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+os.Getenv("STRAPI_JWT"))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
func GetItem(itemId int) (*Item, error) {
	// Get the price id.
	req, _ := http.NewRequest("GET", "http://strapi:1337/api/items/"+strconv.Itoa(itemId), nil)
	req.Header.Set("Authorization", "bearer "+os.Getenv("STRAPI_JWT"))
	client := &http.Client{}
	res, err := client.Do(req)
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
	if item.Data.Attributes.Stock < quantity {
		err = util.NewError("Not enough stock")
		return err
	}
	err = RegisterPrestock(itemId, quantity, item.Data.Attributes.Stock-quantity)
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

	err = RegisterPrestock(itemId, item.Data.Attributes.PreStock-quantity, item.Data.Attributes.Stock+quantity)
	if err != nil {
		return err
	}
	return nil
}

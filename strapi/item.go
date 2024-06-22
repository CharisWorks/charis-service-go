package strapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func RegisterPriceId(itemId int, priceId string) error {
	// Register the price id.
	putRest, _ := json.Marshal(map[string]interface{}{
		"data": map[string]string{
			"price_id": priceId,
		},
	})
	b, _ := io.ReadAll(bytes.NewBuffer(putRest))
	log.Print(string(b))
	log.Print(itemId)
	req, _ := http.NewRequest("PUT", "http://strapi:1337/api/items/"+strconv.Itoa(itemId), bytes.NewBuffer(putRest))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+os.Getenv("STRAPI_JWT"))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	log.Print(res)
	defer res.Body.Close()
	return nil
}

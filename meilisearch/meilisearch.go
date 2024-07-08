package meilisearch

import (
	"encoding/json"

	"github.com/charisworks/charisworks-service-go/util"
	"github.com/meilisearch/meilisearch-go"
)

var Client = meilisearch.NewClient(meilisearch.ClientConfig{
	Host:   util.MEILI_URL,
	APIKey: util.MEILI_MASTER_KEY,
})

func InitMeilisearch() error {
	_, err := Client.Index(util.MEILI_ITEM_INDEX).UpdateDistinctAttribute(util.MEILI_ITEM_INDEX_IDENTIFIER)
	if err != nil {
		return err
	}
	return nil
}
func GetItemByID(itemId string) (*Item, error) {
	var item Item
	err := Client.Index(util.MEILI_ITEM_INDEX).GetDocument(itemId, &meilisearch.DocumentQuery{
		Fields: []string{"*"},
	}, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
func RegisterItemToMeilisearch(item []Item) error {
	var items []map[string]interface{}
	jsonData, err := json.Marshal(item)
	if err != nil {
		return err
	}
	json.Unmarshal(jsonData, &items)

	_, err = Client.Index(util.MEILI_ITEM_INDEX).AddDocuments(items, util.MEILI_ITEM_INDEX_IDENTIFIER)
	if err != nil {
		return err
	}
	return nil
}

func DeleteItemFromMeilisearch(itemId string) error {
	_, err := Client.Index(util.MEILI_ITEM_INDEX).DeleteDocument(itemId)
	if err != nil {
		return err
	}
	return nil
}

func ResetMeilisearch() error {
	_, err := Client.Index(util.MEILI_ITEM_INDEX).DeleteAllDocuments()
	if err != nil {
		return err
	}
	return nil
}

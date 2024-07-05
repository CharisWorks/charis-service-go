package meilisearch

import (
	"encoding/json"

	"github.com/charisworks/charisworks-service-go/util"
	"github.com/meilisearch/meilisearch-go"
)

func RegisterItemToMeilisearch(item []Item) error {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   util.MEILI_URL,
		APIKey: util.MEILI_MASTER_KEY,
	})
	var items []map[string]interface{}
	jsonData, err := json.Marshal(item)
	if err != nil {
		return err
	}
	json.Unmarshal(jsonData, &items)

	_, err = client.Index(util.MEILI_ITEM_INDEX).AddDocuments(items, util.MEILI_ITEM_INDEX_IDENTIFIER)
	if err != nil {
		return err
	}
	return nil
}

func DeleteItemFromMeilisearch(itemId string) error {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   util.MEILI_URL,
		APIKey: util.MEILI_MASTER_KEY,
	})
	_, err := client.Index(util.MEILI_ITEM_INDEX).DeleteDocument(itemId)
	if err != nil {
		return err
	}
	return nil
}
func ResetMeilisearch() error {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   util.MEILI_URL,
		APIKey: util.MEILI_MASTER_KEY,
	})
	_, err := client.Index(util.MEILI_ITEM_INDEX).DeleteAllDocuments()
	if err != nil {
		return err
	}
	return nil
}

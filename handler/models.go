package handler

type transactionPayload struct {
	ItemId   int `json:"item_id"`
	Quantity int `json:"quantity"`
}

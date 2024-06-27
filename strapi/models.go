package strapi

import "time"

type Event struct {
	EventName `json:"event"`
	Model     ModelName `json:"model"`
}
type EventName string
type ItemEvent struct {
	Event
	ItemEntry
}
type ItemEntry struct {
	Entry struct {
		ID          int    `json:"id"`
		Name        string `json:"item_name"`
		Price       int    `json:"price"`
		Stock       int    `json:"stock"`
		Description string `json:"description"`
		Genre       string `json:"genre"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
		PublishedAt string `json:"publishedAt"`
		PriceId     string `json:"price_id"`
	} `json:"entry"`
}
type Item struct {
	Data struct {
		Id         int `json:"id"`
		Attributes struct {
			Name        string `json:"item_name"`
			Price       int    `json:"price"`
			Stock       int    `json:"stock"`
			Description string `json:"description"`
			Genre       string `json:"genre"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
			PublishedAt string `json:"publishedAt"`
			PriceId     string `json:"price_id"`
			PreStock    int    `json:"pre_stock"`
			Worker      struct {
				Data struct {
					ID         int `json:"id"`
					Attributes struct {
						UserName        string `json:"user_name"`
						Description     string `json:"description"`
						CreatedAt       string `json:"createdAt"`
						UpdatedAt       string `json:"updatedAt"`
						PublishedAt     string `json:"publishedAt"`
						StripeAccountID string `json:"stripe_account_id"`
					} `json:"attributes"`
				} `json:"data"`
			} `json:"worker"`
		} `json:"attributes"`
	} `json:"data"`
}
type Data struct {
	ID         int        `json:"id"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	TransactionID string      `json:"transaction_id"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	PublishedAt   time.Time   `json:"publishedAt"`
	Status        string      `json:"status"`
	TrackingID    interface{} `json:"tracking_id"`
	Quantity      int         `json:"quantity"`
	TransferID    interface{} `json:"transfer_id"`
	PaymentIntent string      `json:"payment_intent"`
	Item          Item        `json:"item"`
}

type Transaction struct {
	Data []Data `json:"data"`
	Meta Meta   `json:"meta"`
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	PageCount int `json:"pageCount"`
	Total     int `json:"total"`
}
type CustomerInfo struct {
	Data []struct {
		Id         int `json:"id"`
		Attributes struct {
			CheckoutId string `json:"checkout_id"`
		} `json:"attributes"`
	} `json:"data"`
}

const (
	Create EventName = "entry.create"
	Update EventName = "entry.update"
	Delete EventName = "entry.delete"
)

type ModelName string

const (
	ItemModel ModelName = "item"
)

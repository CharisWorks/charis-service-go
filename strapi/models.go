package strapi

type EventItem struct {
	Event `json:"event"`
	Model string `json:"model"`
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
type Event string
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
		} `json:"attributes"`
	} `json:"data"`
}

const (
	Create Event = "create"
)

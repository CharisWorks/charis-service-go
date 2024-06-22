package strapi

type Item struct {
	Event `json:"event"`
	Model string    `json:"model"`
	Entry ItemEntry `json:"entry"`
}
type Event string

const (
	Create Event = "create"
)

type ItemEntry struct {
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
}

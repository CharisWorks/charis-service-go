package meilisearch

type Item struct {
	ID                string   `json:"id"`
	ItemName          string   `json:"item_name"`
	Price             int      `json:"price"`
	Stock             int      `json:"stock"`
	Description       string   `json:"description"`
	Genre             string   `json:"genre"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`
	PublishedAt       string   `json:"publishedAt"`
	Worker            string   `json:"worker"`
	WorkerDescription string   `json:"worker_description"`
	ThumbnailUrl      string   `json:"thumbnail_url"`
	Images            []Images `json:"images"`
	Tags              []string `json:"tags"`
}

type Images struct {
	SmallUrl  string `json:"small_url"`
	MediumUrl string `json:"medium_url"`
	LargeUrl  string `json:"large_url"`
}

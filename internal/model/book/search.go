package model

type BookSearchResult struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	PublishDate int64    `json:"publish_date"`
	Image       string   `json:"image"`
}

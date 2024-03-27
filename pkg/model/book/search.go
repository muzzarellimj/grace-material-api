package model

type BookSearchResult struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	PublishDate int64  `json:"publish_date"`
	ISBN10      string `json:"isbn_10"`
	ISBN13      string `json:"isbn_13"`
}

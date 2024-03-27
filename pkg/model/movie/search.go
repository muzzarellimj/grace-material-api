package model

type MovieSearchResult struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate int64  `json:"release_date"`
	Image       string `json:"image"`
}

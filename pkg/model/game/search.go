package model

type GameSearchResult struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate int64  `json:"release_date"`
	Image       string `json:"image"`
}

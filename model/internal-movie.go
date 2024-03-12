package model

import "time"

type Movie struct {
	ID                  int       `json:"id"`
	Title               string    `json:"title"`
	Description         *string   `json:"description"`
	Tagline             *string   `json:"tagline"`
	Genres              *[]string `json:"genres"`
	ProductionCompanies *[]string `json:"production_companies"`
	ReleaseDate         time.Time `json:"release_date"`
	Runtime             *int      `json:"runtime"`
	Image               *string   `json:"image"`
	ReferenceImdb       string    `json:"reference_imdb"`
	ReferenceTmdb       string    `json:"reference_tmdb"`
}
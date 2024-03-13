package model

type TMDBMovieDetailResponse struct {
	ID                  int                     `json:"id"`
	Title               string                  `json:"title"`
	Tagline             string                  `json:"tagline"`
	Overview            string                  `json:"overview"`
	Genres              []TMDBGenre             `json:"genres"`
	ProductionCompanies []TMDBProductionCompany `json:"production_companies"`
	ReleaseDate         string                  `json:"release_date"`
	Runtime             int                     `json:"runtime"`
	Image               string                  `json:"poster_path"`
}

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

type TMDBMovieSearchResponse struct {
	Page         int                     `json:"page"`
	Results      []TMDBMovieSearchResult `json:"results"`
	TotalPages   int                     `json:"total_pages"`
	TotalResults int                     `json:"total_results"`
}

type TMDBMovieSearchResult struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Image       string `json:"poster_path"`
}

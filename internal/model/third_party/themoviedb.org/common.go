package model

type TMDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TMDBProductionCompany struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"logo_path"`
}

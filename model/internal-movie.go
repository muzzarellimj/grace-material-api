package model

type Movie struct {
	ID                  int                          `json:"id"`
	Title               string                       `json:"title"`
	Tagline             string                       `json:"tagline"`
	Description         string                       `json:"description"`
	Genres              []MGenreFragment             `json:"genres"`
	ProductionCompanies []MProductionCompanyFragment `json:"production_companies"`
	ReleaseDate         string                       `json:"release_date"`
	Runtime             int                          `json:"runtime"`
	Image               string                       `json:"image"`
	Reference           int                          `json:"reference"`
}

type MovieFragment struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Tagline     string `json:"tagline"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	Runtime     int    `json:"runtime"`
	Image       string `json:"image"`
	Reference   int    `json:"reference"`
}

type MGenreFragment struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Reference int    `json:"reference"`
}

type MProductionCompanyFragment struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Reference int    `json:"reference"`
}

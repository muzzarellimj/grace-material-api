package model

type Movie struct {
	ID                  int                              `json:"id"`
	Title               string                           `json:"title"`
	Tagline             string                           `json:"tagline"`
	Description         string                           `json:"description"`
	Genres              []MovieGenreFragment             `json:"genres"`
	ProductionCompanies []MovieProductionCompanyFragment `json:"production_companies"`
	ReleaseDate         string                           `json:"release_date"`
	Runtime             int                              `json:"runtime"`
	Image               string                           `json:"image"`
	Reference           int                              `json:"reference"`
}

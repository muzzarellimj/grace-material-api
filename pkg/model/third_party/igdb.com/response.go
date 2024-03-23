package model

type IGDBGameResponse struct {
	ID                int                         `json:"id"`
	Title             string                      `json:"name"`
	Summary           string                      `json:"summary"`
	Storyline         string                      `json:"storyline"`
	Franchises        []IGDBNestedNamedResource   `json:"franchises"`
	Genres            []IGDBNestedNamedResource   `json:"genres"`
	InvolvedCompanies []IGDBNestedInvolvedCompany `json:"involved_companies"`
	Platforms         []IGDBNestedNamedResource   `json:"platforms"`
	ReleaseDate       int                         `json:"first_release_date"`
	Cover             IGDBNestedCover             `json:"cover"`
}

type IGDBCompanyResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type IGDBNestedCover struct {
	Hash string `json:"image_id"`
}

type IGDBNestedInvolvedCompany struct {
	Company   int  `json:"company"`
	Developer bool `json:"developer"`
}

type IGDBNestedNamedResource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

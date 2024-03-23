package model

type IGDBGameResponse struct {
	ID                        int    `json:"id"`
	Title                     string `json:"name"`
	Summary                   string `json:"summary"`
	Storyline                 string `json:"storyline"`
	FranchiseReferences       []int  `json:"franchises"`
	GenreReferences           []int  `json:"genres"`
	InvolvedCompanyReferences []int  `json:"involved_companies"`
	PlatformReferences        []int  `json:"platforms"`
	ReleaseDate               int    `json:"first_release_date"`
	CoverReference            int    `json:"cover"`
}

type IGDBInvolvedCompanyResponse struct {
	Company int `json:"company"`
}

type IGDBCompanyResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Common wrapper struct for franchises, genres, and platforms.
type IGDBNamedResourceResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type IGDBCoverResponse struct {
	Hash string `json:"image_id"`
}

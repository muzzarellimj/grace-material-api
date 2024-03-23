package model

type Game struct {
	ID          int                     `json:"id"`
	Title       string                  `json:"title"`
	Summary     string                  `json:"summary"`
	Storyline   string                  `json:"storyline"`
	Franchises  []GameFranchiseFragment `json:"franchises"`
	Genres      []GameGenreFragment     `json:"genres"`
	Platforms   []GamePlatformFragment  `json:"platforms"`
	Studios     []GameStudioFragment    `json:"studios"`
	ReleaseDate int                     `json:"release_date"`
	Image       string                  `json:"image"`
	Reference   int                     `json:"reference"`
}

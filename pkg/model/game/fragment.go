package model

type GameFragment struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Storyline   string `json:"storyline"`
	ReleaseDate string `json:"release_date"`
	Image       string `json:"image"`
	Reference   int    `json:"reference"`
}

type GameFranchiseFragment struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Reference int    `json:"reference"`
}

type GameGenreFragment struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Reference int    `json:"reference"`
}

type GamePlatformFragment struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Reference int    `json:"reference"`
}

type GameStudioFragment struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Reference   int    `json:"reference"`
}

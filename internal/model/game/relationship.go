package model

type GameFranchiseRelationship struct {
	Game      int `json:"game"`
	Franchise int `json:"franchise"`
}

type GameGenreRelationship struct {
	Game  int `json:"game"`
	Genre int `json:"genre"`
}

type GamePlatformRelationship struct {
	Game     int `json:"game"`
	Platform int `json:"platform"`
}

type GameStudioRelationship struct {
	Game   int `json:"game"`
	Studio int `json:"studio"`
}

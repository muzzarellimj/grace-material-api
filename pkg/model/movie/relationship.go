package model

type MovieGenreRelationship struct {
	ID    int    `json:"id"`
	Movie string `json:"movie"`
	Genre string `json:"genre"`
}

type MovieProductionCompanyRelationship struct {
	ID                int    `json:"id"`
	Movie             string `json:"movie"`
	ProductionCompany string `json:"production_company"`
}

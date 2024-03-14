package model

type MovieGenreRelationship struct {
	ID    int `json:"id"`
	Movie int `json:"movie"`
	Genre int `json:"genre"`
}

type MovieProductionCompanyRelationship struct {
	ID                int `json:"id"`
	Movie             int `json:"movie"`
	ProductionCompany int `json:"production_company"`
}

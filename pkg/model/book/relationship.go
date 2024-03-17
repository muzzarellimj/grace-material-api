package model

type BookAuthorRelationship struct {
	Book   int `json:"book"`
	Author int `json:"author"`
}

type BookPublisherRelationship struct {
	Book      int `json:"book"`
	Publisher int `json:"publisher"`
}

type BookTopicRelationship struct {
	Book  int `json:"book"`
	Topic int `json:"topic"`
}

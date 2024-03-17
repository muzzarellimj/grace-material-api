package model

type BookFragment struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	Subtitle         string `json:"subtitle"`
	Description      string `json:"description"`
	PublishDate      string `json:"publish_date"`
	Pages            int    `json:"pages"`
	ISBN10           string `json:"isbn10"`
	ISBN13           string `json:"isbn13"`
	Image            string `json:"image"`
	EditionReference string `json:"edition_reference"`
	WorkReference    string `json:"work_reference"`
}

type BookAuthorFragment struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstname"`
	MiddleName string `json:"middlename"`
	LastName   string `json:"lastname"`
	Biography  string `json:"biography"`
	Image      string `json:"image"`
	Reference  string `json:"reference"`
}

type BookPublisherFragment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BookTopicFragment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

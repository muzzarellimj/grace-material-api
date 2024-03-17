package model

type Book struct {
	ID               int                     `json:"id"`
	Title            string                  `json:"title"`
	Subtitle         string                  `json:"subtitle"`
	Description      string                  `json:"description"`
	Authors          []BookAuthorFragment    `json:"authors"`
	Publishers       []BookPublisherFragment `json:"publishers"`
	Topics           []BookTopicFragment     `json:"topics"`
	PublishDate      string                  `json:"publish_date"`
	Pages            int                     `json:"pages"`
	ISBN10           string                  `json:"isbn10"`
	ISBN13           string                  `json:"isbn13"`
	Image            string                  `json:"image"`
	EditionReference string                  `json:"edition_reference"`
	WorkReference    string                  `json:"work_reference"`
}

package model

type OLAuthorResponse struct {
	ID        string `json:"key"`
	Name      string `json:"name"`
	Biography string `json:"bio"`
	Images    []int  `json:"photos"`
}

type OLEditionResponse struct {
	ID          string                `json:"key"`
	Title       string                `json:"title"`
	Subtitle    string                `json:"subtitle"`
	Authors     []OLResourceReference `json:"authors"`
	Publishers  []string              `json:"publishers"`
	PublishDate string                `json:"publish_date"`
	Format      string                `json:"physical_format"`
	Pages       int                   `json:"number_of_pages"`
	Images      []int                 `json:"covers"`
	ISBN10      []string              `json:"isbn_10"`
	ISBN13      []string              `json:"isbn_13"`
	Works       []OLResourceReference `json:"works"`
}

type OLWorkResponse struct {
	ID          string      `json:"key"`
	Title       string      `json:"title"`
	Description interface{} `json:"description"`
	Subjects    []string    `json:"subjects"`
}

type OLBookSearchResponse struct {
	Results      []OLBookSearchResult `json:"docs"`
	TotalResults int                  `json:"numFound"`
}

type OLBookSearchResult struct {
	ID          []string `json:"edition_key"`
	Title       string   `json:"title"`
	Authors     []string `json:"author_name"`
	PublishDate []string `json:"publish_date"`
}

type OLResourceReference struct {
	ID string `json:"key"`
}

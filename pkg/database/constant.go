package database

// Table names, which may be shared between databases.
const (
	TableBookFragments              = "books"
	TableBookAuthorFragments        = "authors"
	TableBookPublisherFragments     = "publishers"
	TableBookTopicFragments         = "topics"
	TableBookAuthorRelationships    = "books_authors"
	TableBookPublisherRelationships = "books_publishers"
	TableBookTopicRelationships     = "books_topics"

	TableGameFragments              = "games"
	TableGameFranchiseFragments     = "franchises"
	TableGameGenreFragments         = "genres"
	TableGamePlatformFragments      = "platforms"
	TableGameStudioFragments        = "studios"
	TableGameFranchiseRelationships = "games_franchises"
	TableGameGenreRelationships     = "games_genres"
	TableGamePlatformRelationships  = "games_platforms"
	TableGameStudioRelationships    = "games_studios"

	TableMovieFragments                      = "movies"
	TableMovieGenreFragments                 = "genres"
	TableMovieProductionCompanyFragments     = "production_companies"
	TableMovieGenreRelationships             = "movies_genres"
	TableMovieProductionCompanyRelationships = "movies_production_companies"
)

// Properties (or columns names) per database table.
var (
	PropertiesBookFragments              = []string{"id", "title", "subtitle", "description", "publish_date", "pages", "isbn10", "isbn13", "image", "edition_reference", "work_reference"}
	PropertiesBookAuthorFragments        = []string{"id", "firstname", "middlename", "lastname", "biography", "image", "reference"}
	PropertiesBookPublisherFragments     = []string{"id", "name"}
	PropertiesBookTopicFragments         = []string{"id", "name"}
	PropertiesBookAuthorRelationships    = []string{"book", "author"}
	PropertiesBookPublisherRelationships = []string{"book", "publisher"}
	PropertiesBookTopicRelationships     = []string{"book", "topic"}

	PropertiesGameFragments              = []string{"id", "title", "summary", "storyline", "release_date", "image", "reference"}
	PropertiesGameFranchiseFragments     = []string{"id", "name", "reference"}
	PropertiesGameGenreFragments         = []string{"id", "name", "reference"}
	PropertiesGamePlatformFragments      = []string{"id", "name", "reference"}
	PropertiesGameStudioFragments        = []string{"id", "name", "description", "reference"}
	PropertiesGameFranchiseRelationships = []string{"game", "franchise"}
	PropertiesGameGenreRelationships     = []string{"game", "genre"}
	PropertiesGamePlatformRelationships  = []string{"game", "platform"}
	PropertiesGameStudioRelationships    = []string{"game", "studio"}

	PropertiesMovieFragments                      = []string{"title", "tagline", "description", "release_date", "runtime", "image", "reference"}
	PropertiesMovieGenreFragments                 = []string{"name", "reference"}
	PropertiesMovieProductionCompanyFragments     = []string{"name", "image", "reference"}
	PropertiesMovieGenreRelationships             = []string{"movie", "genre"}
	PropertiesMovieProductionCompanyRelationships = []string{"movie", "production_company"}
)

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
	PropertiesBookFragments              = []string{"title", "subtitle", "description", "publish_date", "pages", "isbn10", "isbn13", "image", "edition_reference", "work_reference"}
	PropertiesBookAuthorFragments        = []string{"first_name", "middle_name", "last_name", "biography", "image", "reference"}
	PropertiesBookPublisherFragments     = []string{"name"}
	PropertiesBookTopicFragments         = []string{"name"}
	PropertiesBookAuthorRelationships    = []string{"book", "author"}
	PropertiesBookPublisherRelationships = []string{"book", "publisher"}
	PropertiesBookTopicRelationships     = []string{"book", "topic"}

	PropertiesGameFragments              = []string{"title", "summary", "storyline", "release_date", "image", "reference"}
	PropertiesGameFranchiseFragments     = []string{"name", "reference"}
	PropertiesGameGenreFragments         = []string{"name", "reference"}
	PropertiesGamePlatformFragments      = []string{"name", "reference"}
	PropertiesGameStudioFragments        = []string{"name", "description", "reference"}
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

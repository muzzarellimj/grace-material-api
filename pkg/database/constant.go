package database

// Table names, which may be shared between databases.
const (
	TableMovieFragments                      = "movies"
	TableMovieGenreFragments                 = "genres"
	TableMovieProductionCompanyFragments     = "production_companies"
	TableMovieGenreRelationships             = "movies_genres"
	TableMovieProductionCompanyRelationships = "movies_production_companies"
)

// Properties (or columns names) per database table.
var (
	PropertiesMovieFragments                      = []string{"title", "tagline", "description", "release_date", "runtime", "image", "reference"}
	PropertiesMovieGenreFragments                 = []string{"name", "reference"}
	PropertiesMovieProductionCompanyFragments     = []string{"name", "image", "reference"}
	PropertiesMovieGenreRelationships             = []string{"movie", "genre"}
	PropertiesMovieProductionCompanyRelationships = []string{"movie", "production_company"}
)

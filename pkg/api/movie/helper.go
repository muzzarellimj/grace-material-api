package api

import (
	"fmt"
	"os"

	db "github.com/muzzarellimj/grace-material-api/database"
	database "github.com/muzzarellimj/grace-material-api/pkg/database"
	"github.com/muzzarellimj/grace-material-api/pkg/database/service"
	model "github.com/muzzarellimj/grace-material-api/pkg/model/movie"
)

func fetchMovie(field string, value string) (model.Movie, string, error) {
	var connection db.PgxConnection = db.MovieConnection

	const errorResponseMessage string = "Unable to fetch movie metadata and map to supported data structure."

	movieFragment, err := service.FetchFragment[model.MovieFragment](connection, database.TableMovies, fmt.Sprintf("%s=%s", field, value))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch movie: %v\n", err)

		return model.Movie{}, errorResponseMessage, err
	}

	movieGenreRelationships, err := service.FetchRelationshipSlice[model.MovieGenreRelationship](connection, database.TableMoviesGenres, fmt.Sprintf("movie=%d", movieFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between movies and genres: %v\n", err)

		return model.Movie{}, errorResponseMessage, err
	}

	var genreFragments []model.MovieGenreFragment

	for _, relationship := range movieGenreRelationships {
		genreFragment, err := service.FetchFragment[model.MovieGenreFragment](connection, database.TableGenres, fmt.Sprintf("id=%d", relationship.Genre))

		if err == nil {
			genreFragments = append(genreFragments, genreFragment)
		}
	}

	movieProductionCompanyRelationships, err := service.FetchRelationshipSlice[model.MovieProductionCompanyRelationship](connection, database.TableMoviesProductionCompanies, fmt.Sprintf("movie=%d", movieFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between movies and production companies: %v\n", err)

		return model.Movie{}, errorResponseMessage, err
	}

	var productionCompanyFragments []model.MovieProductionCompanyFragment

	for _, relationship := range movieProductionCompanyRelationships {
		productionCompanyFragment, err := service.FetchFragment[model.MovieProductionCompanyFragment](connection, database.TableProductionCompanies, fmt.Sprintf("id=%d", relationship.ProductionCompany))

		if err == nil {
			productionCompanyFragments = append(productionCompanyFragments, productionCompanyFragment)
		}
	}

	return mapMovie(movieFragment, genreFragments, productionCompanyFragments), "", nil
}

func mapMovie(movieFragment model.MovieFragment, genreFragments []model.MovieGenreFragment, productionCompanyFragments []model.MovieProductionCompanyFragment) model.Movie {
	return model.Movie{
		ID:                  movieFragment.ID,
		Title:               movieFragment.Title,
		Tagline:             movieFragment.Tagline,
		Description:         movieFragment.Description,
		Genres:              genreFragments,
		ProductionCompanies: productionCompanyFragments,
		ReleaseDate:         movieFragment.ReleaseDate,
		Runtime:             movieFragment.Runtime,
		Image:               movieFragment.Image,
		Reference:           movieFragment.Reference,
	}
}

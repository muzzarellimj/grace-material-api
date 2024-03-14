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
	movieFragment, err := fetchMovieFragment(field, value)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch movie: %v\n", err)

		return model.Movie{}, "Unable to fetch and map movie metadata.", err
	}

	movieGenreRelationships, err := fetchMovieGenreRelationships(movieFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between movies and genres: %v\n", err)

		return model.Movie{}, "Unable to fetch and map movie metadata.", err
	}

	var genreFragments []model.MovieGenreFragment

	for _, relationship := range movieGenreRelationships {
		genreFragment, err := fetchMovieGenreFragment("id", fmt.Sprint(relationship.Genre))

		if err == nil {
			genreFragments = append(genreFragments, genreFragment)
		}
	}

	movieProductionCompanyRelationships, err := fetchMovieProductionCompanyRelationships(movieFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between movies and production companies: %v\n", err)

		return model.Movie{}, "Unable to fetch and map movie metadata.", err
	}

	var productionCompanyFragments []model.MovieProductionCompanyFragment

	for _, relationship := range movieProductionCompanyRelationships {
		productionCompanyFragment, err := fetchMovieProductionCompanyFragment("id", fmt.Sprint(relationship.ProductionCompany))

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

func fetchMovieFragment(field string, value string) (model.MovieFragment, error) {
	response, err := service.FetchFragmentSlice[model.MovieFragment](db.MovieConnection, database.TableMovies, fmt.Sprintf("%s=%s", field, value))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch movie fragment slice: %v\n", err)

		return model.MovieFragment{}, err
	}

	if len(response) > 0 {
		return response[0], nil
	}

	return model.MovieFragment{}, nil
}

func fetchMovieGenreFragment(field string, value string) (model.MovieGenreFragment, error) {
	response, err := service.FetchFragmentSlice[model.MovieGenreFragment](db.MovieConnection, database.TableGenres, fmt.Sprintf("%s=%s", field, value))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch genre fragment slice: %v\n", err)

		return model.MovieGenreFragment{}, err
	}

	if len(response) > 0 {
		return response[0], nil
	}

	return model.MovieGenreFragment{}, nil
}

func fetchMovieProductionCompanyFragment(field string, value string) (model.MovieProductionCompanyFragment, error) {
	response, err := service.FetchFragmentSlice[model.MovieProductionCompanyFragment](db.MovieConnection, database.TableProductionCompanies, fmt.Sprintf("%s=%s", field, value))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch production company fragment slice: %v\n", err)

		return model.MovieProductionCompanyFragment{}, err
	}

	if len(response) > 0 {
		return response[0], nil
	}

	return model.MovieProductionCompanyFragment{}, nil
}

func fetchMovieGenreRelationships(movieFragment model.MovieFragment) ([]model.MovieGenreRelationship, error) {
	response, err := service.FetchFragmentSlice[model.MovieGenreRelationship](db.MovieConnection, database.TableMoviesGenres, fmt.Sprintf("movie=%d", movieFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between movies and genres: %v\n", err)

		return []model.MovieGenreRelationship{}, err
	}

	return response, nil
}

func fetchMovieProductionCompanyRelationships(movieFragment model.MovieFragment) ([]model.MovieProductionCompanyRelationship, error) {
	response, err := service.FetchFragmentSlice[model.MovieProductionCompanyRelationship](db.MovieConnection, database.TableMoviesProductionCompanies, fmt.Sprintf("movie=%d", movieFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between movies and production companies: %v\n", err)

		return []model.MovieProductionCompanyRelationship{}, err
	}

	return response, nil
}

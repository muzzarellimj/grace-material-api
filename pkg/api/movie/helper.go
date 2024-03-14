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
	movieFragment, message, err := fetchMovieFragment(field, value)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch movie: %v\n", err)

		return model.Movie{}, message, err
	}

	return mapMovie(movieFragment), "", nil
}

func mapMovie(movieFragment model.MovieFragment) model.Movie {
	return model.Movie{
		ID:                  movieFragment.ID,
		Title:               movieFragment.Title,
		Tagline:             movieFragment.Tagline,
		Description:         movieFragment.Description,
		Genres:              []model.MovieGenreFragment{},
		ProductionCompanies: []model.MovieProductionCompanyFragment{},
		ReleaseDate:         movieFragment.ReleaseDate,
		Runtime:             movieFragment.Runtime,
		Image:               movieFragment.Image,
		Reference:           movieFragment.Reference,
	}
}

func fetchMovieFragment(field string, value string) (model.MovieFragment, string, error) {
	response, err := service.FetchFragmentSlice[model.MovieFragment](db.MovieConnection, database.TableMovies, fmt.Sprintf("%s=%s", field, value))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch movie fragment slice: %v\n", err)

		return model.MovieFragment{}, "Unable to fetch necessary movie data.", err
	}

	if len(response) > 0 {
		return response[0], "", nil
	}

	return model.MovieFragment{}, "", nil
}

package helper

import (
	"fmt"
	"os"

	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/movie"
)

func FetchMovie(constraint string) (model.Movie, error) {
	zero := model.Movie{}

	movieFragment, err := service.FetchFragment[model.MovieFragment](database.Connection, database.TableMovieFragments, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch movie with constraint '%s': %v\n", constraint, err)

		return zero, err
	}

	genreFragmentSlice, err := fetchGenreFragmentSlice(movieFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch genres related to movie '%d': %v\n", movieFragment.ID, err)
	}

	productionCompanyFragmentSlice, err := fetchProductionCompanyFragmentSlice(movieFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch production companies related to movie '%d': %v\n", movieFragment.ID, err)
	}

	movie := mapMovie(movieFragment, genreFragmentSlice, productionCompanyFragmentSlice)

	return movie, nil
}

func FetchMovieSlice(constraintSlice []string) ([]model.Movie, []error) {
	var movieSlice []model.Movie
	var errorSlice []error

	for _, constraint := range constraintSlice {
		movie, err := FetchMovie(constraint)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch and map movie with constraint '%s': %v\n", constraint, err)

			errorSlice = append(errorSlice, err)
		}

		if movie.ID != 0 {
			movieSlice = append(movieSlice, movie)
		}
	}

	return movieSlice, errorSlice
}

func fetchGenreFragmentSlice(movieFragment model.MovieFragment) ([]model.MovieGenreFragment, error) {
	zero := []model.MovieGenreFragment{}

	movieGenreRelationshipSlice, err := service.FetchRelationshipSlice[model.MovieGenreRelationship](database.Connection, database.TableMovieGenreRelationships, fmt.Sprintf("movie=%d", movieFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between movie '%d' and genres: %v\n", movieFragment.ID, err)

		return zero, err
	}

	var genreFragmentSlice []model.MovieGenreFragment

	for _, relationship := range movieGenreRelationshipSlice {
		genreFragment, err := service.FetchFragment[model.MovieGenreFragment](database.Connection, database.TableMovieGenreFragments, fmt.Sprintf("id=%d", relationship.Genre))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch genre '%d': %v\n", relationship.Genre, err)
		}

		if genreFragment.ID != 0 {
			genreFragmentSlice = append(genreFragmentSlice, genreFragment)
		}
	}

	return genreFragmentSlice, nil
}

func fetchProductionCompanyFragmentSlice(movieFragment model.MovieFragment) ([]model.MovieProductionCompanyFragment, error) {
	zero := []model.MovieProductionCompanyFragment{}

	movieProductionCompanyRelationshipSlice, err := service.FetchRelationshipSlice[model.MovieProductionCompanyRelationship](database.Connection, database.TableMovieProductionCompanyRelationships, fmt.Sprintf("movie=%d", movieFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between movie '%d' and production companies: %v\n", movieFragment.ID, err)

		return zero, err
	}

	var productionCompanyFragmentSlice []model.MovieProductionCompanyFragment

	for _, relationship := range movieProductionCompanyRelationshipSlice {
		productionCompanyFragment, err := service.FetchFragment[model.MovieProductionCompanyFragment](database.Connection, database.TableMovieProductionCompanyFragments, fmt.Sprintf("id=%d", relationship.ProductionCompany))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch production company '%d': %v\n", relationship.ProductionCompany, err)
		}

		if productionCompanyFragment.ID != 0 {
			productionCompanyFragmentSlice = append(productionCompanyFragmentSlice, productionCompanyFragment)
		}
	}

	return productionCompanyFragmentSlice, nil
}

func mapMovie(movieFragment model.MovieFragment, genreFragmentSlice []model.MovieGenreFragment, productionCompanyFragmentSlice []model.MovieProductionCompanyFragment) model.Movie {
	if genreFragmentSlice == nil {
		genreFragmentSlice = make([]model.MovieGenreFragment, 0)
	}

	if productionCompanyFragmentSlice == nil {
		productionCompanyFragmentSlice = make([]model.MovieProductionCompanyFragment, 0)
	}

	return model.Movie{
		ID:                  movieFragment.ID,
		Title:               movieFragment.Title,
		Tagline:             movieFragment.Tagline,
		Description:         movieFragment.Description,
		Genres:              genreFragmentSlice,
		ProductionCompanies: productionCompanyFragmentSlice,
		ReleaseDate:         movieFragment.ReleaseDate,
		Runtime:             movieFragment.Runtime,
		Image:               movieFragment.Image,
		Reference:           movieFragment.Reference,
	}
}

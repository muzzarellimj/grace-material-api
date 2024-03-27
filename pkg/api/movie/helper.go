package api

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/muzzarellimj/grace-material-api/pkg/api/movie/helper"
	"github.com/muzzarellimj/grace-material-api/pkg/database"
	"github.com/muzzarellimj/grace-material-api/pkg/database/connection"
	"github.com/muzzarellimj/grace-material-api/pkg/database/service"
	model "github.com/muzzarellimj/grace-material-api/pkg/model/movie"
	tmodel "github.com/muzzarellimj/grace-material-api/pkg/model/third_party/themoviedb.org"
	"github.com/muzzarellimj/grace-material-api/pkg/util"
)

func fetchMovie(constraint string) (model.Movie, string, error) {
	var connection connection.PgxPool = connection.Movie

	const errorResponseMessage string = "Unable to fetch movie metadata and map to supported data structure."

	movieFragment, err := service.FetchFragment[model.MovieFragment](connection, database.TableMovieFragments, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch movie: %v\n", err)

		return model.Movie{}, errorResponseMessage, err
	}

	movieGenreRelationships, err := service.FetchRelationshipSlice[model.MovieGenreRelationship](connection, database.TableMovieGenreRelationships, fmt.Sprintf("movie=%d", movieFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between movies and genres: %v\n", err)

		return model.Movie{}, errorResponseMessage, err
	}

	var genreFragments []model.MovieGenreFragment

	for _, relationship := range movieGenreRelationships {
		genreFragment, err := service.FetchFragment[model.MovieGenreFragment](connection, database.TableMovieGenreFragments, fmt.Sprintf("id=%d", relationship.Genre))

		if err == nil {
			genreFragments = append(genreFragments, genreFragment)
		}
	}

	movieProductionCompanyRelationships, err := service.FetchRelationshipSlice[model.MovieProductionCompanyRelationship](connection, database.TableMovieProductionCompanyRelationships, fmt.Sprintf("movie=%d", movieFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between movies and production companies: %v\n", err)

		return model.Movie{}, errorResponseMessage, err
	}

	var productionCompanyFragments []model.MovieProductionCompanyFragment

	for _, relationship := range movieProductionCompanyRelationships {
		productionCompanyFragment, err := service.FetchFragment[model.MovieProductionCompanyFragment](connection, database.TableMovieProductionCompanyFragments, fmt.Sprintf("id=%d", relationship.ProductionCompany))

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

func storeMovie(tmdbMovie tmodel.TMDBMovieDetailResponse) (int, error) {
	movieId, err := storeMovieFragment(tmdbMovie)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store movie: %v\n", err)

		return 0, err
	}

	genreIdSlice := storeGenreFragments(tmdbMovie.Genres)
	productionCompanyIdSlice := storeProductionCompanyFragments(tmdbMovie.ProductionCompanies)

	storeMovieGenreRelationships(movieId, genreIdSlice)
	storeMovieProductionCompanyRelationships(movieId, productionCompanyIdSlice)

	return movieId, nil
}

func storeMovieFragment(tmdbMovie tmodel.TMDBMovieDetailResponse) (int, error) {
	movieId, err := service.StoreFragment(connection.Movie, database.TableMovieFragments, database.PropertiesMovieFragments, pgx.NamedArgs{
		"title":        tmdbMovie.Title,
		"tagline":      tmdbMovie.Tagline,
		"description":  tmdbMovie.Overview,
		"release_date": util.ParseDateTime(tmdbMovie.ReleaseDate),
		"runtime":      tmdbMovie.Runtime,
		"image":        helper.FormatImagePath(tmdbMovie.Image),
		"reference":    tmdbMovie.ID,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store movie fragment: %v\n", err)

		return 0, err
	}

	return movieId, nil
}

func storeGenreFragments(genres []tmodel.TMDBGenre) []int {
	var genreIdSlice []int

	for _, genre := range genres {
		existingGenreFragment, err := service.FetchFragment[model.MovieGenreFragment](connection.Movie, database.TableMovieGenreFragments, fmt.Sprintf("reference=%d", genre.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing genre fragment: %v\n", err)

			continue
		}

		if existingGenreFragment.ID != 0 {
			genreIdSlice = append(genreIdSlice, existingGenreFragment.ID)

			continue
		}

		storedGenreId, err := service.StoreFragment(connection.Movie, database.TableMovieGenreFragments, database.PropertiesMovieGenreFragments, pgx.NamedArgs{
			"name":      genre.Name,
			"reference": genre.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new genre fragment: %v\n", err)

			continue
		}

		if storedGenreId != 0 {
			genreIdSlice = append(genreIdSlice, storedGenreId)
		}
	}

	return genreIdSlice
}

func storeProductionCompanyFragments(productionCompanies []tmodel.TMDBProductionCompany) []int {
	var productionCompanyIdSlice []int

	for _, productionCompany := range productionCompanies {
		existingProductionCompanyFragment, err := service.FetchFragment[model.MovieProductionCompanyFragment](connection.Movie, database.TableMovieProductionCompanyFragments, fmt.Sprintf("reference=%d", productionCompany.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing production company fragment: %v\n", err)

			continue
		}

		if existingProductionCompanyFragment.ID != 0 {
			productionCompanyIdSlice = append(productionCompanyIdSlice, existingProductionCompanyFragment.ID)

			continue
		}

		storedProductionCompanyId, err := service.StoreFragment(connection.Movie, database.TableMovieProductionCompanyFragments, database.PropertiesMovieProductionCompanyFragments, pgx.NamedArgs{
			"name":      productionCompany.Name,
			"image":     helper.FormatImagePath(productionCompany.Image),
			"reference": productionCompany.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new production company fragment: %v\n", err)

			continue
		}

		if storedProductionCompanyId != 0 {
			productionCompanyIdSlice = append(productionCompanyIdSlice, storedProductionCompanyId)
		}
	}

	return productionCompanyIdSlice
}

func storeMovieGenreRelationships(movieId int, genreIdSlice []int) {
	for _, genreId := range genreIdSlice {
		err := service.StoreRelationship(connection.Movie, database.TableMovieGenreRelationships, database.PropertiesMovieGenreRelationships, pgx.NamedArgs{
			"movie": movieId,
			"genre": genreId,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new relationship between movie '%d' and genre '%d': %v\n", movieId, genreId, err)
		}
	}
}

func storeMovieProductionCompanyRelationships(movieId int, productionCompanyIdSlice []int) {
	for _, productionCompanyId := range productionCompanyIdSlice {
		err := service.StoreRelationship(connection.Movie, database.TableMovieProductionCompanyRelationships, database.PropertiesMovieProductionCompanyRelationships, pgx.NamedArgs{
			"movie":              movieId,
			"production_company": productionCompanyId,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new relationship between movie '%d' and production company '%d': %v\n", movieId, productionCompanyId, err)
		}
	}
}

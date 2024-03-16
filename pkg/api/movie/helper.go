package api

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/muzzarellimj/grace-material-api/pkg/database"
	"github.com/muzzarellimj/grace-material-api/pkg/database/connection"
	"github.com/muzzarellimj/grace-material-api/pkg/database/service"
	model "github.com/muzzarellimj/grace-material-api/pkg/model/movie"
	tmodel "github.com/muzzarellimj/grace-material-api/pkg/model/third_party/themoviedb.org"
)

func fetchMovie(field string, value string) (model.Movie, string, error) {
	var connection connection.PgxPool = connection.Movie

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

func storeMovie(tmdbMovie tmodel.TMDBMovieDetailResponse) (int, error) {
	var connection connection.PgxPool = connection.Movie

	storedMovieId, err := service.StoreFragment(connection, database.TableMovies, database.PropertiesMovies, pgx.NamedArgs{
		"title":        tmdbMovie.Title,
		"tagline":      tmdbMovie.Tagline,
		"description":  tmdbMovie.Overview,
		"release_date": tmdbMovie.ReleaseDate,
		"runtime":      tmdbMovie.Runtime,
		"image":        tmdbMovie.Image,
		"reference":    tmdbMovie.ID,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store movie fragment: %v\n", err)

		return 0, err
	}

	var storedGenreIds []int

	for _, genre := range tmdbMovie.Genres {
		genreFragment, err := service.FetchFragment[model.MovieGenreFragment](connection, database.TableGenres, fmt.Sprintf("reference=%d", genre.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch genre fragment: %v\n", err)
		}

		if genreFragment.ID != 0 {
			storedGenreIds = append(storedGenreIds, genreFragment.ID)

			continue
		}

		storedGenreId, err := service.StoreFragment(connection, database.TableGenres, database.PropertiesMoviesGenres, pgx.NamedArgs{
			"name":      genre.Name,
			"reference": genre.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store genre fragment: %v\n", err)

			continue
		}

		if storedGenreId != 0 {
			storedGenreIds = append(storedGenreIds, storedGenreId)
		}
	}

	var storedProductionCompanyIds []int

	for _, productionCompany := range tmdbMovie.ProductionCompanies {
		productionCompanyFragment, err := service.FetchFragment[model.MovieProductionCompanyFragment](connection, database.TableProductionCompanies, fmt.Sprintf("reference=%d", productionCompany.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch production company fragment: %v\n", err)
		}

		if productionCompanyFragment.ID != 0 {
			storedProductionCompanyIds = append(storedProductionCompanyIds, productionCompanyFragment.ID)

			continue
		}

		storedProductionCompanyId, err := service.StoreFragment(connection, database.TableProductionCompanies, database.PropertiesMoviesProductionCompanies, pgx.NamedArgs{
			"name":      productionCompany.Name,
			"image":     productionCompany.Image,
			"reference": productionCompany.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store genre fragment: %v\n", err)

			continue
		}

		if storedProductionCompanyId != 0 {
			storedProductionCompanyIds = append(storedProductionCompanyIds, storedProductionCompanyId)
		}
	}

	for _, storedGenreId := range storedGenreIds {
		err := service.StoreRelationship(connection, database.TableMoviesGenres, database.PropertiesMoviesGenresRelationships, pgx.NamedArgs{
			"movie": storedMovieId,
			"genre": storedGenreId,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store relationship between movie and genre: %v\n", err)
		}
	}

	for _, storedProductionCompanyId := range storedProductionCompanyIds {
		err := service.StoreRelationship(connection, database.TableMoviesProductionCompanies, database.PropertiesMoviesProductionCompaniesRelationships, pgx.NamedArgs{
			"movie":              storedMovieId,
			"production_company": storedProductionCompanyId,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store relationship between movie and production company: %v\n", err)
		}
	}

	return storedMovieId, nil
}

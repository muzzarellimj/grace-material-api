package helper

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/muzzarellimj/grace-material-api/internal/api/game/helper"
	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/movie"
	TMDBModel "github.com/muzzarellimj/grace-material-api/internal/model/third_party/themoviedb.org"
	"github.com/muzzarellimj/grace-material-api/internal/util"
)

func ProcessMovieStorage(movie TMDBModel.TMDBMovieDetailResponse) (int, error) {
	movieId, err := storeMovieFragment(movie)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store movie '%d' : %v\n", movie.ID, err)

		return 0, err
	}

	genreIdSlice := processGenreFragmentSlice(movie.Genres)
	productionCompanyIdSlice := processProductionCompanyFragmentSlice(movie.ProductionCompanies)

	service.StoreRelationshipSlice(database.Connection, database.TableMovieGenreRelationships, database.PropertiesMovieGenreRelationships, service.RelationshipSliceArgument{
		SourceName:          "movie",
		SourceArgument:      movieId,
		DestinationName:     "genre",
		DestinationArgument: genreIdSlice,
	})

	service.StoreRelationshipSlice(database.Connection, database.TableMovieProductionCompanyRelationships, database.PropertiesMovieProductionCompanyRelationships, service.RelationshipSliceArgument{
		SourceName:          "movie",
		SourceArgument:      movieId,
		DestinationName:     "production_company",
		DestinationArgument: productionCompanyIdSlice,
	})

	return movieId, nil
}

func storeMovieFragment(movie TMDBModel.TMDBMovieDetailResponse) (int, error) {
	movieId, err := service.StoreFragment(database.Connection, database.TableMovieFragments, database.PropertiesMovieFragments, pgx.NamedArgs{
		"title":        movie.Title,
		"tagline":      movie.Tagline,
		"description":  movie.Overview,
		"release_date": util.ParseDateTime(movie.ReleaseDate),
		"runtime":      movie.Runtime,
		"image":        helper.FormatImagePath(movie.Image),
		"reference":    movie.ID,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store movie '%d' fragment: %v\n", movie.ID, err)

		return 0, err
	}

	return movieId, nil
}

func processGenreFragmentSlice(genres []TMDBModel.TMDBGenre) []int {
	var genreIdSlice []int

	for _, genre := range genres {
		existingGenreFragment, err := service.FetchFragment[model.MovieGenreFragment](database.Connection, database.TableMovieGenreFragments, fmt.Sprintf("reference=%d", genre.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing genre '%d' fragment: %v\n", genre.ID, err)

			continue
		}

		if existingGenreFragment.ID != 0 {
			genreIdSlice = append(genreIdSlice, existingGenreFragment.ID)

			continue
		}

		genreId, err := service.StoreFragment(database.Connection, database.TableMovieGenreFragments, database.PropertiesMovieGenreFragments, pgx.NamedArgs{
			"name":      genre.Name,
			"reference": genre.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new genre '%d' fragment: %v\n", genre.ID, err)
		}

		if genreId != 0 {
			genreIdSlice = append(genreIdSlice, genreId)
		}
	}

	return genreIdSlice
}

func processProductionCompanyFragmentSlice(productionCompanies []TMDBModel.TMDBProductionCompany) []int {
	var productionCompanyIdSlice []int

	for _, productionCompany := range productionCompanies {
		existingProductionCompanyFragment, err := service.FetchFragment[model.MovieProductionCompanyFragment](database.Connection, database.TableMovieProductionCompanyFragments, fmt.Sprintf("reference=%d", productionCompany.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing production company '%d' fragment: %v\n", productionCompany.ID, err)

			continue
		}

		if existingProductionCompanyFragment.ID != 0 {
			productionCompanyIdSlice = append(productionCompanyIdSlice, existingProductionCompanyFragment.ID)

			continue
		}

		productionCompanyId, err := service.StoreFragment(database.Connection, database.TableMovieProductionCompanyFragments, database.PropertiesMovieProductionCompanyFragments, pgx.NamedArgs{
			"name":      productionCompany.Name,
			"image":     productionCompany.Image,
			"reference": productionCompany.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new production company '%d' fragment: %v\n", productionCompany.ID, err)
		}

		if productionCompanyId != 0 {
			productionCompanyIdSlice = append(productionCompanyIdSlice, productionCompanyId)
		}
	}

	return productionCompanyIdSlice
}

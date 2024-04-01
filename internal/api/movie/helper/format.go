package helper

import (
	"fmt"

	model "github.com/muzzarellimj/grace-material-api/internal/model/movie"
	TMDBModel "github.com/muzzarellimj/grace-material-api/internal/model/third_party/themoviedb.org"
	"github.com/muzzarellimj/grace-material-api/internal/util"
)

func MapSearchResultSlice(input []TMDBModel.TMDBMovieSearchResult) []model.MovieSearchResult {
	var resultSlice []model.MovieSearchResult

	for _, result := range input {
		mappedResult := model.MovieSearchResult{
			ID:          result.ID,
			Title:       result.Title,
			ReleaseDate: util.ParseDateTime(result.ReleaseDate),
			Image:       FormatImagePath(result.Image),
		}

		resultSlice = append(resultSlice, mappedResult)
	}

	return resultSlice
}

func FormatImagePath(path string) string {
	return fmt.Sprint("https://image.tmdb.org/t/p/original", path)
}

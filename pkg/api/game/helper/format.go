package helper

import (
	"fmt"

	model "github.com/muzzarellimj/grace-material-api/pkg/model/game"
	IGDBModel "github.com/muzzarellimj/grace-material-api/pkg/model/third_party/igdb.com"
)

func MapSearchResultSlice(input []IGDBModel.IGDBGameSearchResponse) []model.GameSearchResult {
	var resultSlice []model.GameSearchResult

	for _, result := range input {
		mappedResult := model.GameSearchResult{
			ID:          result.ID,
			Title:       result.Title,
			ReleaseDate: int64(result.ReleaseDate),
			Image:       FormatImagePath(result.Cover.Hash),
		}

		resultSlice = append(resultSlice, mappedResult)
	}

	return resultSlice
}

func FormatImagePath(hash string) string {
	return fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_%s/%s.jpg", "cover_big", hash)
}

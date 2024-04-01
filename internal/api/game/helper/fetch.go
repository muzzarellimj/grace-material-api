package helper

import (
	"fmt"
	"os"

	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/connection"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/game"
)

func FetchGame(constraint string) (model.Game, error) {
	gameFragment, err := service.FetchFragment[model.GameFragment](connection.Game, database.TableGameFragments, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch game with constraint '%s': %v\n", constraint, err)

		return model.Game{}, err
	}

	franchiseFragmentSlice, err := fetchFranchiseFragmentSlice(gameFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch franchises related to game '%d': %v\n", gameFragment.ID, err)
	}

	genreFragmentSlice, err := fetchGenreFragmentSlice(gameFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch genres related to game '%d': %v\n", gameFragment.ID, err)
	}

	platformFragmentSlice, err := fetchPlatformFragmentSlice(gameFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch platforms related to game '%d': %v\n", gameFragment.ID, err)
	}

	studioFragmentSlice, err := fetchStudioFragmentSlice(gameFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch studios related to game '%d': %v\n", gameFragment.ID, err)
	}

	return mapGame(gameFragment, franchiseFragmentSlice, genreFragmentSlice, platformFragmentSlice, studioFragmentSlice), nil
}

func fetchFranchiseFragmentSlice(gameFragment model.GameFragment) ([]model.GameFranchiseFragment, error) {
	gameFranchiseRelationshipSlice, err := service.FetchRelationshipSlice[model.GameFranchiseRelationship](connection.Game, database.TableGameFranchiseRelationships, fmt.Sprintf("game=%d", gameFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between game '%d' and franchises: %v\n", gameFragment.ID, err)

		return []model.GameFranchiseFragment{}, err
	}

	var franchiseFragmentSlice []model.GameFranchiseFragment

	for _, relationship := range gameFranchiseRelationshipSlice {
		franchiseFragment, err := service.FetchFragment[model.GameFranchiseFragment](connection.Game, database.TableGameFranchiseFragments, fmt.Sprintf("id=%d", relationship.Franchise))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch franchise '%d': %v\n", relationship.Franchise, err)
		}

		if franchiseFragment.ID != 0 {
			franchiseFragmentSlice = append(franchiseFragmentSlice, franchiseFragment)
		}
	}

	return franchiseFragmentSlice, nil
}

func fetchGenreFragmentSlice(gameFragment model.GameFragment) ([]model.GameGenreFragment, error) {
	gameGenreRelationshipSlice, err := service.FetchRelationshipSlice[model.GameGenreRelationship](connection.Game, database.TableGameGenreRelationships, fmt.Sprintf("game=%d", gameFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between game '%d' and genres: %v\n", gameFragment.ID, err)

		return []model.GameGenreFragment{}, err
	}

	var genreFragmentSlice []model.GameGenreFragment

	for _, relationship := range gameGenreRelationshipSlice {
		genreFragment, err := service.FetchFragment[model.GameGenreFragment](connection.Game, database.TableGameGenreFragments, fmt.Sprintf("id=%d", relationship.Genre))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch genre '%d': %v\n", relationship.Genre, err)
		}

		if genreFragment.ID != 0 {
			genreFragmentSlice = append(genreFragmentSlice, genreFragment)
		}
	}

	return genreFragmentSlice, nil
}

func fetchPlatformFragmentSlice(gameFragment model.GameFragment) ([]model.GamePlatformFragment, error) {
	gamePlatformRelationshipSlice, err := service.FetchRelationshipSlice[model.GamePlatformRelationship](connection.Game, database.TableGamePlatformRelationships, fmt.Sprintf("game=%d", gameFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between game '%d' and platforms: %v\n", gameFragment.ID, err)

		return []model.GamePlatformFragment{}, err
	}

	var platformFragmentSlice []model.GamePlatformFragment

	for _, relationship := range gamePlatformRelationshipSlice {
		platformFragment, err := service.FetchFragment[model.GamePlatformFragment](connection.Game, database.TableGamePlatformFragments, fmt.Sprintf("id=%d", relationship.Platform))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch platform '%d': %v\n", relationship.Platform, err)
		}

		if platformFragment.ID != 0 {
			platformFragmentSlice = append(platformFragmentSlice, platformFragment)
		}
	}

	return platformFragmentSlice, nil
}

func fetchStudioFragmentSlice(gameFragment model.GameFragment) ([]model.GameStudioFragment, error) {
	gameStudioRelationshipSlice, err := service.FetchRelationshipSlice[model.GameStudioRelationship](connection.Game, database.TableGameStudioRelationships, fmt.Sprintf("game=%d", gameFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between game '%d' and studios: %v\n", gameFragment.ID, err)

		return []model.GameStudioFragment{}, err
	}

	var studioFragmentSlice []model.GameStudioFragment

	for _, relationship := range gameStudioRelationshipSlice {
		studioFragment, err := service.FetchFragment[model.GameStudioFragment](connection.Game, database.TableGameStudioFragments, fmt.Sprintf("id=%d", relationship.Studio))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch studio '%d': %v\n", relationship.Studio, err)
		}

		if studioFragment.ID != 0 {
			studioFragmentSlice = append(studioFragmentSlice, studioFragment)
		}
	}

	return studioFragmentSlice, nil
}

func mapGame(gameFragment model.GameFragment, franchiseFragmentSlice []model.GameFranchiseFragment, genreFragmentSlice []model.GameGenreFragment, platformFragmentSlice []model.GamePlatformFragment, studioFragmentSlice []model.GameStudioFragment) model.Game {
	return model.Game{
		ID:          gameFragment.ID,
		Title:       gameFragment.Title,
		Summary:     gameFragment.Summary,
		Storyline:   gameFragment.Storyline,
		Franchises:  franchiseFragmentSlice,
		Genres:      genreFragmentSlice,
		Platforms:   platformFragmentSlice,
		Studios:     studioFragmentSlice,
		ReleaseDate: gameFragment.ReleaseDate,
		Image:       gameFragment.Image,
		Reference:   gameFragment.Reference,
	}
}

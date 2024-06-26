package helper

import (
	"fmt"
	"os"

	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/game"
)

func FetchGame(constraint string) (model.Game, error) {
	zero := model.Game{}

	gameFragment, err := service.FetchFragment[model.GameFragment](database.Connection, database.TableGameFragments, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch game with constraint '%s': %v\n", constraint, err)

		return zero, err
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

	game := mapGame(gameFragment, franchiseFragmentSlice, genreFragmentSlice, platformFragmentSlice, studioFragmentSlice)

	return game, nil
}

func FetchGameSlice(constraintSlice []string) ([]model.Game, []error) {
	var gameSlice []model.Game
	var errSlice []error

	for _, constraint := range constraintSlice {
		game, err := FetchGame(constraint)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch and map game with constraint '%s': %v\n", constraint, err)

			errSlice = append(errSlice, err)
		}

		if game.ID != 0 {
			gameSlice = append(gameSlice, game)
		}
	}

	return gameSlice, errSlice
}

func FetchGameExistenceSlice() ([]int, []error) {
	idSlice, err := service.FetchExistenceSlice(database.Connection, database.TableGameFragments)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch existence slice: %v\n", err)

		return []int{}, []error{err}
	}

	if len(idSlice) == 0 {
		fmt.Fprint(os.Stderr, "Existence slice appears to be empty.")

		return []int{}, nil
	}

	return idSlice, nil
}

func fetchFranchiseFragmentSlice(gameFragment model.GameFragment) ([]model.GameFranchiseFragment, error) {
	gameFranchiseRelationshipSlice, err := service.FetchRelationshipSlice[model.GameFranchiseRelationship](database.Connection, database.TableGameFranchiseRelationships, fmt.Sprintf("game=%d", gameFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between game '%d' and franchises: %v\n", gameFragment.ID, err)

		return []model.GameFranchiseFragment{}, err
	}

	var franchiseFragmentSlice []model.GameFranchiseFragment

	for _, relationship := range gameFranchiseRelationshipSlice {
		franchiseFragment, err := service.FetchFragment[model.GameFranchiseFragment](database.Connection, database.TableGameFranchiseFragments, fmt.Sprintf("id=%d", relationship.Franchise))

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
	gameGenreRelationshipSlice, err := service.FetchRelationshipSlice[model.GameGenreRelationship](database.Connection, database.TableGameGenreRelationships, fmt.Sprintf("game=%d", gameFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between game '%d' and genres: %v\n", gameFragment.ID, err)

		return []model.GameGenreFragment{}, err
	}

	var genreFragmentSlice []model.GameGenreFragment

	for _, relationship := range gameGenreRelationshipSlice {
		genreFragment, err := service.FetchFragment[model.GameGenreFragment](database.Connection, database.TableGameGenreFragments, fmt.Sprintf("id=%d", relationship.Genre))

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
	gamePlatformRelationshipSlice, err := service.FetchRelationshipSlice[model.GamePlatformRelationship](database.Connection, database.TableGamePlatformRelationships, fmt.Sprintf("game=%d", gameFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between game '%d' and platforms: %v\n", gameFragment.ID, err)

		return []model.GamePlatformFragment{}, err
	}

	var platformFragmentSlice []model.GamePlatformFragment

	for _, relationship := range gamePlatformRelationshipSlice {
		platformFragment, err := service.FetchFragment[model.GamePlatformFragment](database.Connection, database.TableGamePlatformFragments, fmt.Sprintf("id=%d", relationship.Platform))

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
	gameStudioRelationshipSlice, err := service.FetchRelationshipSlice[model.GameStudioRelationship](database.Connection, database.TableGameStudioRelationships, fmt.Sprintf("game=%d", gameFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between game '%d' and studios: %v\n", gameFragment.ID, err)

		return []model.GameStudioFragment{}, err
	}

	var studioFragmentSlice []model.GameStudioFragment

	for _, relationship := range gameStudioRelationshipSlice {
		studioFragment, err := service.FetchFragment[model.GameStudioFragment](database.Connection, database.TableGameStudioFragments, fmt.Sprintf("id=%d", relationship.Studio))

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
	if franchiseFragmentSlice == nil {
		franchiseFragmentSlice = make([]model.GameFranchiseFragment, 0)
	}

	if genreFragmentSlice == nil {
		genreFragmentSlice = make([]model.GameGenreFragment, 0)
	}

	if platformFragmentSlice == nil {
		platformFragmentSlice = make([]model.GamePlatformFragment, 0)
	}

	if studioFragmentSlice == nil {
		studioFragmentSlice = make([]model.GameStudioFragment, 0)
	}

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

package helper

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	IGDBAPI "github.com/muzzarellimj/grace-material-api/pkg/api/third_party/igdb.com"
	"github.com/muzzarellimj/grace-material-api/pkg/database"
	"github.com/muzzarellimj/grace-material-api/pkg/database/connection"
	"github.com/muzzarellimj/grace-material-api/pkg/database/service"
	model "github.com/muzzarellimj/grace-material-api/pkg/model/game"
	IGDBModel "github.com/muzzarellimj/grace-material-api/pkg/model/third_party/igdb.com"
)

func ProcessGameStorage(game IGDBModel.IGDBGameResponse) (int, error) {
	gameId, err := storeGameFragment(game)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store game '%d' fragment: %v\n", game.ID, err)

		return 0, err
	}

	franchiseIdSlice := processFranchiseFragmentSlice(game.Franchises)
	genreIdSlice := processGenreFragmentSlice(game.Genres)
	platformIdSlice := processPlatformFragmentSlice(game.Platforms)
	studioIdSlice := processStudioFragmentSlice(game.InvolvedCompanies)

	service.StoreRelationshipSlice(connection.Game, database.TableGameFranchiseRelationships, database.PropertiesGameFranchiseRelationships, service.RelationshipSliceArgument{
		SourceName:          "game",
		SourceArgument:      gameId,
		DestinationName:     "franchise",
		DestinationArgument: franchiseIdSlice,
	})

	service.StoreRelationshipSlice(connection.Game, database.TableGameGenreRelationships, database.PropertiesGameGenreRelationships, service.RelationshipSliceArgument{
		SourceName:          "game",
		SourceArgument:      gameId,
		DestinationName:     "genre",
		DestinationArgument: genreIdSlice,
	})

	service.StoreRelationshipSlice(connection.Game, database.TableGamePlatformRelationships, database.PropertiesGamePlatformRelationships, service.RelationshipSliceArgument{
		SourceName:          "game",
		SourceArgument:      gameId,
		DestinationName:     "platform",
		DestinationArgument: platformIdSlice,
	})

	service.StoreRelationshipSlice(connection.Game, database.TableGameStudioRelationships, database.PropertiesGameStudioRelationships, service.RelationshipSliceArgument{
		SourceName:          "game",
		SourceArgument:      gameId,
		DestinationName:     "studio",
		DestinationArgument: studioIdSlice,
	})

	return gameId, nil
}

func storeGameFragment(game IGDBModel.IGDBGameResponse) (int, error) {
	gameId, err := service.StoreFragment(connection.Game, database.TableGameFragments, database.PropertiesGameFragments, pgx.NamedArgs{
		"title":        game.Title,
		"summary":      game.Summary,
		"storyline":    game.Storyline,
		"release_date": game.ReleaseDate,
		"image":        FormatImagePath(game.Cover.Hash),
		"reference":    game.ID,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to store game '%d' fragment: %v\n", game.ID, err)

		return 0, err
	}

	return gameId, nil
}

func processFranchiseFragmentSlice(franchises []IGDBModel.IGDBNestedNamedResource) []int {
	var franchiseIdSlice []int

	for _, resource := range franchises {
		existingFranchiseFragment, err := service.FetchFragment[model.GameFranchiseFragment](connection.Game, database.TableGameFranchiseFragments, fmt.Sprintf("reference=%d", resource.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing franchise '%d' fragment: %v\n", resource.ID, err)

			continue
		}

		if existingFranchiseFragment.ID != 0 {
			franchiseIdSlice = append(franchiseIdSlice, existingFranchiseFragment.ID)

			continue
		}

		franchiseId, err := service.StoreFragment(connection.Game, database.TableGameFranchiseFragments, database.PropertiesGameFranchiseFragments, pgx.NamedArgs{
			"name":      resource.Name,
			"reference": resource.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new franchise '%d' fragment: %v\n", resource.ID, err)
		}

		if franchiseId != 0 {
			franchiseIdSlice = append(franchiseIdSlice, franchiseId)
		}
	}

	return franchiseIdSlice
}

func processGenreFragmentSlice(genres []IGDBModel.IGDBNestedNamedResource) []int {
	var genreIdSlice []int

	for _, resource := range genres {
		existingGenreFragment, err := service.FetchFragment[model.GameGenreFragment](connection.Game, database.TableGameGenreFragments, fmt.Sprintf("reference=%d", resource.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing genre '%d' fragment: %v\n", resource.ID, err)

			continue
		}

		if existingGenreFragment.ID != 0 {
			genreIdSlice = append(genreIdSlice, existingGenreFragment.ID)

			continue
		}

		genreId, err := service.StoreFragment(connection.Game, database.TableGameGenreFragments, database.PropertiesGameGenreFragments, pgx.NamedArgs{
			"name":      resource.Name,
			"reference": resource.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new genre '%d' fragment: %v\n", resource.ID, err)
		}

		if genreId != 0 {
			genreIdSlice = append(genreIdSlice, genreId)
		}
	}

	return genreIdSlice
}

func processPlatformFragmentSlice(platforms []IGDBModel.IGDBNestedNamedResource) []int {
	var platformIdSlice []int

	for _, resource := range platforms {
		existingPlatformFragment, err := service.FetchFragment[model.GamePlatformFragment](connection.Game, database.TableGamePlatformFragments, fmt.Sprintf("reference=%d", resource.ID))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing platform '%d' fragment: %v\n", resource.ID, err)

			continue
		}

		if existingPlatformFragment.ID != 0 {
			platformIdSlice = append(platformIdSlice, existingPlatformFragment.ID)

			continue
		}

		platformId, err := service.StoreFragment(connection.Game, database.TableGamePlatformFragments, database.PropertiesGamePlatformFragments, pgx.NamedArgs{
			"name":      resource.Name,
			"reference": resource.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new platform '%d' fragment: %v\n", resource.ID, err)
		}

		if platformId != 0 {
			platformIdSlice = append(platformIdSlice, platformId)
		}
	}

	return platformIdSlice
}

func processStudioFragmentSlice(companies []IGDBModel.IGDBNestedInvolvedCompany) []int {
	var studioIdSlice []int

	for _, company := range companies {
		if !company.Developer {
			continue
		}

		existingStudioFragment, err := service.FetchFragment[model.GameStudioFragment](connection.Game, database.TableGameStudioFragments, fmt.Sprintf("reference=%d", company.Company))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch existing studio '%d' fragment: %v\n", company.Company, err)

			continue
		}

		if existingStudioFragment.ID != 0 {
			studioIdSlice = append(studioIdSlice, existingStudioFragment.ID)

			continue
		}

		studio, err := IGDBAPI.IGDBGetResource[IGDBModel.IGDBCompanyResponse](IGDBAPI.IGDBEndpointCompany, fmt.Sprintf("fields id,name,description; where id=%d;", company.Company))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch company '%d' IGDB record: %v\n", company.Company, err)

			continue
		}

		studioId, err := service.StoreFragment(connection.Game, database.TableGameStudioFragments, database.PropertiesGameStudioFragments, pgx.NamedArgs{
			"name":        studio.Name,
			"description": studio.Description,
			"reference":   studio.ID,
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to store new studio '%d' fragment: %v\n", studio.ID, err)
		}

		if studioId != 0 {
			studioIdSlice = append(studioIdSlice, studioId)
		}
	}

	return studioIdSlice
}

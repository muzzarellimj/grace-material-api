package helper

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/game"
)

func UpdateGameFragment(game model.GameFragment) (int, error) {
	id, err := service.UpdateFragment(database.Connection, database.TableGameFragments, database.PropertiesGameFragments, fmt.Sprintf("id=%d", game.ID), pgx.NamedArgs{
		"title":        game.Title,
		"summary":      game.Summary,
		"storyline":    game.Storyline,
		"release_date": game.ReleaseDate,
		"image":        game.Image,
		"reference":    game.Reference,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update game '%d' fragment: %v\n", game.ID, err)

		return 0, err
	}

	return id, nil
}

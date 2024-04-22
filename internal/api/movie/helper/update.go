package helper

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/movie"
)

func UpdateMovieFragment(movie model.MovieFragment) (int, error) {
	id, err := service.UpdateFragment(database.Connection, database.TableMovieFragments, database.PropertiesMovieFragments, fmt.Sprintf("id=%d", movie.ID), pgx.NamedArgs{
		"title":        movie.Title,
		"tagline":      movie.Tagline,
		"description":  movie.Description,
		"release_date": movie.ReleaseDate,
		"runtime":      movie.Runtime,
		"image":        movie.Image,
		"reference":    movie.Reference,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update movie '%d' fragment: %v\n", movie.ID, err)

		return 0, err
	}

	return id, nil
}

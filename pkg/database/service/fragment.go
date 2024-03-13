package service

import (
	"fmt"
	"os"

	"github.com/muzzarellimj/grace-material-api/database"
)

const (
	TableMovies              = "movies"
	TableGenres              = "genres"
	TableProductionCompanies = "production_companies"
)

func FetchFragmentSlice[M interface{}](connection database.PgxConnection, table string, constraint string) ([]M, error) {
	statement, err := database.CreateQuery("*", table, constraint, "")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create fragment slice selection statement: %v\n", err)

		return []M{}, err
	}

	rows, err := database.ExecuteQuery(connection, statement)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute fragment slice selection statement '%s': %v\n", statement, err)

		return []M{}, err
	}

	response, err := database.MapResponse[M](rows)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to map fragment slice selection response: %v\n", err)

		return []M{}, err
	}

	return response, nil
}

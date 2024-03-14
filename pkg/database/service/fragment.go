package service

import (
	"fmt"
	"os"

	"github.com/muzzarellimj/grace-material-api/database"
)

func FetchFragment[M interface{}](connection database.PgxConnection, table string, constraint string) (M, error) {
	var zero M

	fragmentSlice, err := FetchFragmentSlice[M](connection, table, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch initial fragment slice: %v\n", err)

		return zero, err
	}

	if len(fragmentSlice) > 0 {
		return fragmentSlice[0], nil
	}

	return zero, nil
}

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

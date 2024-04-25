package service

import (
	"fmt"
	"os"

	"github.com/muzzarellimj/grace-material-api/internal/database"
)

func FetchExistenceSlice(connection database.PgxPool, table string) ([]int, error) {
	var zero []int

	statement, err := database.CreateQuery("id", table, "", "")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create existence slice selection statement: %v\n", err)

		return zero, err
	}

	rows, err := database.ExecuteQuery(connection, statement)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute existence slice selection statement: %v\n", err)

		return zero, err
	}

	response, err := database.MapQueryResponse[int](rows)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to map existence slice selection response: %v\n", err)

		return zero, err
	}

	return response, nil
}

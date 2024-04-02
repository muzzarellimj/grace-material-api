package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/muzzarellimj/grace-material-api/internal/database"
)

func FetchFragment[M interface{}](connection database.PgxPool, table string, constraint string) (M, error) {
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

func FetchFragmentSlice[M interface{}](connection database.PgxPool, table string, constraint string) ([]M, error) {
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

	response, err := database.MapQueryResponse[M](rows)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to map fragment slice selection response: %v\n", err)

		return []M{}, err
	}

	return response, nil
}

// Store a fragment in the provided table with the provided properties (column names) and named arguments.
//
// Return: the numeric identifier for the stored fragment and nil with success, or 0 and error without.
func StoreFragment(connection database.PgxPool, table string, properties []string, arguments pgx.NamedArgs) (int, error) {
	var names []string

	for _, property := range properties {
		names = append(names, fmt.Sprint("@", property))
	}

	statement := fmt.Sprintf("INSERT INTO %s (%v) VALUES (%v) RETURNING id", table, strings.Join(properties, ","), strings.Join(names, ","))

	tx, err := connection.Begin(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to begin transaction to store fragment: %v\n", err)

		return 0, err
	}

	defer func() {
		err = tx.Rollback(context.Background())

		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			fmt.Fprintf(os.Stderr, "Unable to rollback transaction: %v\n", err)
		}
	}()

	var id int

	err = tx.QueryRow(context.Background(), statement, arguments).Scan(&id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute fragment insertion statement: %v\n", err)

		return 0, err
	}

	err = tx.Commit(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to commit fragment insertion transaction: %v\n", err)

		return 0, err
	}

	return id, nil
}

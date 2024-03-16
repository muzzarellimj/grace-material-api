package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/muzzarellimj/grace-material-api/pkg/database/connection"
)

func FetchRelationship[M interface{}](connection connection.PgxPool, table string, constraint string) (M, error) {
	var zero M

	relationshipSlice, err := FetchFragmentSlice[M](connection, table, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch initial relationship slice: %s\n", err)

		return zero, err
	}

	if len(relationshipSlice) > 0 {
		return relationshipSlice[0], nil
	}

	return zero, nil
}

func FetchRelationshipSlice[M interface{}](connection connection.PgxPool, table string, constraint string) ([]M, error) {
	var zero []M

	relationshipSlice, err := FetchFragmentSlice[M](connection, table, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationship slice: %v\n", err)

		return zero, err
	}

	return relationshipSlice, nil
}

// Store a relationship in the provided table with the provided properties (column names) and named arguments.
//
// Return: nil with success, and error without.
func StoreRelationship(connection connection.PgxPool, table string, properties []string, arguments pgx.NamedArgs) error {
	var names []string

	for _, property := range properties {
		names = append(names, fmt.Sprint("@", property))
	}

	statement := fmt.Sprintf("INSERT INTO %s (%v) VALUES (%v)", table, strings.Join(properties, ","), strings.Join(names, ","))

	tx, err := connection.Begin(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to begin transaction to store relationship: %v\n", err)

		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), statement, arguments)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute relationship insertion statement: %v\n", err)

		return err
	}

	err = tx.Commit(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to commit relationship insertion transaction: %v\n", err)

		return err
	}

	return nil
}

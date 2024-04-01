package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/muzzarellimj/grace-material-api/internal/database/connection"
)

// Create a PostgreSQL query statement with given selection, from, where, and group statements,
// as well as optional directives (e.g., join statements).
//
// Return: built statement and nil with success, empty string and error without.
func CreateQuery(selection string, from string, where string, group string, directives ...string) (string, error) {
	var builder strings.Builder

	if selection == "" || from == "" {
		err := errors.New("unable to without 'selection' and 'from' args")

		fmt.Fprintf(os.Stderr, "Unable to create query without 'selection' and 'from' statements: %v\n", err)

		return "", err
	}

	builder.WriteString(fmt.Sprintf("SELECT %s ", selection))
	builder.WriteString(fmt.Sprintf("FROM %s ", from))

	for _, directive := range directives {
		builder.WriteString(fmt.Sprintf("%s ", directive))
	}

	if where != "" {
		builder.WriteString(fmt.Sprintf("WHERE %s ", where))
	}

	if group != "" {
		builder.WriteString(fmt.Sprintf("GROUP BY %s", group))
	}

	return strings.TrimSpace(builder.String()), nil
}

// Execute a PostgreSQL query within the given database connection and with the given
// query statement.
//
// Return: pgx.Rows-type response and nil with success, nil and error without.
func ExecuteQuery(connection connection.PgxPool, statement string) (pgx.Rows, error) {
	if statement == "" {
		err := errors.New("unable to execute query without 'statement' arg")

		fmt.Fprintf(os.Stderr, "Unable to execute query without 'statement' argument: %v\n", err)

		return nil, err
	}

	response, err := connection.Query(context.Background(), statement)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute query: %v\n", err)

		return nil, err
	}

	return response, nil
}

// Map a PostgreSQL query response to a supported data model slice.
//
// Return: parsed model slice and nil with success, empty model slice and error without.
func MapQueryResponse[M interface{}](rows pgx.Rows) ([]M, error) {
	var response []M

	err := pgxscan.ScanAll(&response, rows)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to map query response to supported data model: %v\n", err)

		return []M{}, err
	}

	return response, nil
}

package database

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Create a PostgreSQL query statement with given selection, from, where, and group statements,
// as well as optional directives (e.g., join statements).
//
// Return: built statement and nil with success, empty string and error without.
func CreateQuery(selection string, from string, where string, group string, directives ...string) (string, error) {
	var builder strings.Builder

	if selection == "" || from == "" {
		fmt.Fprintf(os.Stderr, "Unable to create query without 'selection' and 'from' statements.")

		return "", errors.New("query: unable to create query without selection and from args")
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

	return builder.String(), nil
}

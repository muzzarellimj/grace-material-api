package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool PgxPool

// A wrapper to mask pgxpool.Pool as a local interface.
type PgxPool interface {
	Begin(context context.Context) (pgx.Tx, error)
	Close()
	Query(context context.Context, swl string, args ...any) (pgx.Rows, error)
}

// Connect the Grace database pool and persist the connections as accessible variables.
//
// Return: error with error occurrence, nil without.
func Connect(url string) error {
	fmt.Fprint(os.Stdout, "Connect to Grace database pool...\n")

	connection, err := pgxpool.New(context.Background(), url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to Grace database pool: %v\n", err)

		return err
	}

	if connection == nil {
		fmt.Fprint(os.Stderr, "Unable to persist connection to Grace database pool.")

		return err
	}

	Pool = connection

	return nil
}

// Close the Grace database pool connection.
func Disconnect() {
	Pool.Close()
}

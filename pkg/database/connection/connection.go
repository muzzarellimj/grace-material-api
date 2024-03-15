package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Movie PgxPool

// A wrapper to mask pgxpool.Pool as a local interface.
type PgxPool interface {
	Begin(context context.Context) (pgx.Tx, error)
	Close()
	Query(context context.Context, swl string, args ...any) (pgx.Rows, error)
}

// Connect the Grace database pools and persist the connections as accessible variables.
func Connect(username string, password string, host string, port string) {
	fmt.Fprint(os.Stdout, "Connecting the Grace database pools...")

	Movie = connect(username, password, host, port, os.Getenv("MOVIE_DATABASE_NAME"))
}

// Connect to a named database pool.
//
// Return: database pool in PgxPool wrapper interface.
func connect(username string, password string, host string, port string, database string) *pgxpool.Pool {
	url := fmt.Sprint("postgres://", username, ":", password, "@", host, ":", port, "/", database)

	connection, err := pgxpool.New(context.Background(), url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database pool '%s': %v\n", database, err)

		os.Exit(1)
	}

	return connection
}

// Disconnect the Grace database pools.
func Disconnect() {
	Movie.Close()
}

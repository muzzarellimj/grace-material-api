package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A wrapper interface to mask pgxpool.Pool and control local access properties.
type PgxConnection interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close()
}

var MovieConnection PgxConnection

// Connect to the Grace database pool and persist the connections as exported variables.
func Connect(username string, password string, host string, port string) {
	fmt.Println("Connecting to Grace database pools.")

	MovieConnection = connect(username, password, host, port, os.Getenv("MOVIE_DATABASE_NAME"))
}

// Connect to a named database pool.
//
// Return: database pool in PgxConnection wrapper.
func connect(username string, password string, host string, port string, database string) *pgxpool.Pool {
	url := fmt.Sprint("postgres://", username, ":", password, "@", host, ":", port, "/", database)

	connection, err := pgxpool.New(context.Background(), url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database pool %s: %v\n", database, err)
		os.Exit(1)
	}

	return connection
}

// Disconnect from the Grace database pools.
func Disconnect() {
	MovieConnection.Close()
}

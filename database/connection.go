package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var MovieConnection *pgx.Conn

// Connect to the Grace database pool and persist the connections as exported variables.
func Connect(username string, password string, host string, port string) {
	fmt.Println("Connecting to Grace database pool.")

	MovieConnection = connect(username, password, host, port, os.Getenv("MOVIE_DATABASE_NAME"))
}

// Connect to a named individual database.
//
// Return: database connection.
func connect(username string, password string, host string, port string, database string) *pgx.Conn {
	url := fmt.Sprint("postgres://", username, ":", password, "@", host, ":", port, "/", database)

	connection, err := pgx.Connect(context.Background(), url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %s: %v\n", database, err)
		os.Exit(1)
	}

	return connection
}

// Disconnect from the Grace database pool.
func Disconnect() {
	disconnect(MovieConnection)
}

// Disconnect from a provided database connection.
func disconnect(connection *pgx.Conn) {
	err := connection.Close(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to close connection to a database: %v\n", err)
	}
}

package database_test

import (
	"testing"

	"github.com/muzzarellimj/grace-material-api/database"
	"github.com/muzzarellimj/grace-material-api/model"
	"github.com/pashagolub/pgxmock/v3"
)

func TestCreateQueryReturnsSelectFrom(t *testing.T) {
	expected := "SELECT id FROM example "

	query, err := database.CreateQuery("id", "example", "", "")

	if query != expected || err != nil {
		t.Fatalf("Actual query statement '%s' does not match expected query statement '%s': %v\n", query, expected, err)
	}
}

func TestCreateQueryReturnsSelectFromWhere(t *testing.T) {
	expected := "SELECT id FROM example WHERE id=1 "

	query, err := database.CreateQuery("id", "example", "id=1", "")

	if query != expected || err != nil {
		t.Fatalf("Actual query statement '%s' does not match expected query statement '%s': %v\n", query, expected, err)
	}
}

func TestCreateQueryReturnsSelectFromWhereGroup(t *testing.T) {
	expected := "SELECT id FROM example WHERE id=1 GROUP BY id"

	query, err := database.CreateQuery("id", "example", "id=1", "id")

	if query != expected || err != nil {
		t.Fatalf("Actual query statement '%s' does not match expected query statement '%s': %v\n", query, expected, err)
	}
}

func TestCreateQueryReturnsSelectFromWhereGroupDirective(t *testing.T) {
	expected := "SELECT example.id FROM example JOIN other ON example.id=other.id WHERE id=1 GROUP BY example.id"

	query, err := database.CreateQuery("example.id", "example", "id=1", "example.id", "JOIN other ON example.id=other.id")

	if query != expected || err != nil {
		t.Fatalf("Actual query statement '%s' does not match expected query statement '%s': %v\n", query, expected, err)
	}
}

func TestExecuteQueryReturnsSuccess(t *testing.T) {
	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Unable to create mock database connection pool: %v\n", err)
	}

	defer mock.Close()

	mock.ExpectQuery("SELECT id FROM example WHERE id=1").WillReturnRows(pgxmock.NewRows([]string{"1"}))

	statement, _ := database.CreateQuery("id", "example", "id=1", "")
	_, err = database.ExecuteQuery(mock, statement)

	if err != nil {
		t.Fatalf("Query statement '%s' was unable to execute with success: %v\n", statement, err)
	}

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatalf("Query statement '%s' execution was unable to meet mock expectations: %v\n", statement, err)
	}
}

func TestMapResponseReturnsSlice(t *testing.T) {
	expectedTitle := "Encanto"
	expectedReference := 812

	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Unable to create mock database connection pool: %v\n", err)
	}

	defer mock.Close()

	mock.ExpectQuery("SELECT x FROM movies WHERE id=1").WillReturnRows(pgxmock.
		NewRows([]string{"id", "title", "tagline", "description", "genres", "production_companies", "release_date", "runtime", "image", "reference"}).
		AddRow(1, "Encanto", nil, nil, nil, nil, "1992-11-25", nil, nil, 812))

	statement, _ := database.CreateQuery("x", "movies", "id=1", "")
	rows, _ := database.ExecuteQuery(mock, statement)
	response, err := database.MapResponse[model.Movie](rows)

	if err != nil {
		t.Fatalf("Unable to map response '%v' to internal movie struct: %v\n", response, err)
	}

	if len(response) < 1 {
		t.Fatalf("Unable to map response '%v' to internal movie struct, but without without error.", response)
	}

	if response[0].Title != expectedTitle || response[0].Reference != expectedReference {
		t.Fatalf("Actual title '%s' and reference identifier '%d' do not match expected title '%s' and reference identifier '%d'.", response[0].Title, response[0].Reference, expectedTitle, expectedReference)
	}
}

func TestMapResponseReturnsEmptySlice(t *testing.T) {
	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Unable to create mock database connection pool: %v\n", err)
	}

	defer mock.Close()

	mock.ExpectQuery("SELECT x FROM movies WHERE id=1").WillReturnRows(pgxmock.
		NewRows([]string{"id", "title", "description", "tagline", "genres", "production_companies", "release_date", "runtime", "image", "reference_imdb", "reference_tmdb"}))

	statement, _ := database.CreateQuery("x", "movies", "id=1", "")
	rows, _ := database.ExecuteQuery(mock, statement)
	response, err := database.MapResponse[model.Movie](rows)

	if err != nil {
		t.Fatalf("Unable to map response '%v' to empty internal movie struct slice: %v\n", response, err)
	}

	if len(response) > 0 {
		t.Fatalf("Unable to map response '%v' to empty internal movie struct slice, but without without error.", response)
	}
}

func TestMapResponseReturnsTypeMismatchError(t *testing.T) {
	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Unable to create mock database connection pool: %v\n", err)
	}

	defer mock.Close()

	mock.ExpectQuery("SELECT x FROM movies WHERE id=1").WillReturnRows(pgxmock.
		NewRows([]string{"id", "title", "description", "tagline", "genres", "production_companies", "release_date", "runtime", "image", "reference_imdb", "reference_tmdb"}).
		AddRow(1, "Encanto", nil, nil, nil, nil, "incorrect input type", nil, nil, "tt0103639", "812"))

	statement, _ := database.CreateQuery("x", "movies", "id=1", "")
	rows, _ := database.ExecuteQuery(mock, statement)
	response, err := database.MapResponse[model.Movie](rows)

	if err == nil {
		t.Fatalf("Unable to report error while mapping response '%v' to internal movie struct slice.", response)
	}
}

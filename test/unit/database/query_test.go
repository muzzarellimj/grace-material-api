package database_test

import (
	"testing"
	"time"

	"github.com/muzzarellimj/grace-material-api/database"
	"github.com/muzzarellimj/grace-material-api/model"
	"github.com/pashagolub/pgxmock/v3"
)

func TestCreateQuerySelectionFrom(t *testing.T) {
	expected := "SELECT id FROM example "

	query, err := database.CreateQuery("id", "example", "", "")

	if query != expected || err != nil {
		t.Fatalf("Expected query statement %s does not match actual query statement %s: %v\n", expected, query, err)
	}
}

func TestCreateQuerySelectionFromWhere(t *testing.T) {
	expected := "SELECT id FROM example WHERE id=1 "

	query, err := database.CreateQuery("id", "example", "id=1", "")

	if query != expected || err != nil {
		t.Fatalf("Expected query statement %s does not match actual query statement %s: %v\n", expected, query, err)
	}
}

func TestCreateQuerySelectionFromWhereGroup(t *testing.T) {
	expected := "SELECT id FROM example WHERE id=1 GROUP BY id"

	query, err := database.CreateQuery("id", "example", "id=1", "id")

	if query != expected || err != nil {
		t.Fatalf("Expected query statement %s does not match actual query statement %s: %v\n", expected, query, err)
	}
}

func TestCreateQuerySelectionFromWhereGroupDirective(t *testing.T) {
	expected := "SELECT example.id FROM example JOIN other ON example.id=other.id WHERE id=1 GROUP BY example.id"

	query, err := database.CreateQuery("example.id", "example", "id=1", "example.id", "JOIN other ON example.id=other.id")

	if query != expected || err != nil {
		t.Fatalf("Expected query statement %s does not match actual query statement %s: %v\n", expected, query, err)
	}
}

func TestExecuteQuerySelects(t *testing.T) {
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
	expectedImdbId := "tt0103639"

	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Unable to create mock database connection pool: %v\n", err)
	}

	defer mock.Close()

	mock.ExpectQuery("SELECT x FROM movies WHERE id=1").WillReturnRows(pgxmock.
		NewRows([]string{"id", "title", "description", "tagline", "genres", "production_companies", "release_date", "runtime", "image", "reference_imdb", "reference_tmdb"}).
		AddRow(1, "Encanto", nil, nil, nil, nil, time.Time{}, nil, nil, "tt0103639", "812"))

	statement, _ := database.CreateQuery("x", "movies", "id=1", "")
	rows, _ := database.ExecuteQuery(mock, statement)
	response, err := database.MapResponse[model.Movie](rows)

	if err != nil {
		t.Fatalf("Unable to map response '%v' to internal movie struct: %v\n", response, err)
	}

	if len(response) < 1 {
		t.Fatalf("Unable to map response '%v' to internal movie struct, but without without error.", response)
	}

	if response[0].Title != expectedTitle || response[0].ReferenceImdb != expectedImdbId {
		t.Fatalf("Actual title '%s' and IMDB ID '%s' do not match expected title '%s' and IMDB ID '%s'.", response[0].Title, response[0].ReferenceImdb, expectedTitle, expectedImdbId)
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

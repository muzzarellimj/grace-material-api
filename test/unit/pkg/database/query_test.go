package database_test

import (
	"testing"

	"github.com/muzzarellimj/grace-material-api/pkg/database"
	model "github.com/muzzarellimj/grace-material-api/pkg/model/movie"
	"github.com/pashagolub/pgxmock/v3"
)

func TestCreateQueryReturnsSelectFrom(t *testing.T) {
	expected := "SELECT id FROM movies"

	query, err := database.CreateQuery("id", "movies", "", "")

	if err != nil {
		t.Fatalf("Unable to create test query statement: %v\n", err)
	}

	if query != expected {
		t.Fatalf("Actual query statement '%s' does not match expected query statement '%s'.", query, expected)
	}
}

func TestCreateQueryReturnsSelectFromWhere(t *testing.T) {
	expected := "SELECT id FROM genres WHERE id=1"

	query, err := database.CreateQuery("id", "genres", "id=1", "")

	if err != nil {
		t.Fatalf("Unable to create test query statement: %v\n", err)
	}

	if query != expected {
		t.Fatalf("Actual query statement '%s' does not match expected query statement '%s'.", query, expected)
	}
}

func TestCreateQueryReturnsSelectFromWhereGroupBy(t *testing.T) {
	expected := "SELECT id FROM genres WHERE id=1 GROUP BY id"

	query, err := database.CreateQuery("id", "genres", "id=1", "id")

	if err != nil {
		t.Fatalf("Unable to create test query statement: %v\n", err)
	}

	if query != expected {
		t.Fatalf("Actual query statement '%s' does not match expected query statement '%s'.", query, expected)
	}
}

func TestCreateQueryReturnsSelectFromWhereGroupByDirective(t *testing.T) {
	expected := "SELECT m.id FROM movies m JOIN movies_production_companies p ON m.id=p.movie WHERE m.id=1 GROUP BY m.id"

	query, err := database.CreateQuery("m.id", "movies m", "m.id=1", "m.id", "JOIN movies_production_companies p ON m.id=p.movie")

	if err != nil {
		t.Fatalf("Unable to create test query statement: %v\n", err)
	}

	if query != expected {
		t.Fatalf("Actual query statement '%s' does not match expected query statement '%s'.", query, expected)
	}
}

func TestCreateQueryHandlesEmptySelectionArg(t *testing.T) {
	query, err := database.CreateQuery("", "production_companies", "", "")

	if err == nil {
		t.Fatal("Unable to catch error with missing 'selection' argument.")
	}

	if query != "" {
		t.Fatalf("Actual query statement '%s' does not match expected empty query statement.", query)
	}
}

func TestCreateQueryHandlesEmptyFromArg(t *testing.T) {
	query, err := database.CreateQuery("id", "", "", "")

	if err == nil {
		t.Fatal("Unable to catch error with missing 'from' argument.")
	}

	if query != "" {
		t.Fatalf("Actual query statement '%s' does not match expected empty query statement.", query)
	}
}

func TestExecuteQueryReturnsRows(t *testing.T) {
	mock := createMockConnection(t)

	defer mock.Close()

	mock.ExpectQuery("SELECT id FROM genres WHERE id=1").
		WillReturnRows(pgxmock.
			NewRows([]string{"id", "name", "reference"}).
			AddRow(1, "Action", 0))

	statement, _ := database.CreateQuery("id", "genres", "id=1", "")
	_, err := database.ExecuteQuery(mock, statement)

	if err != nil {
		t.Fatalf("Unable to execute query statement '%s': %v\n", statement, err)
	}

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatalf("Unable to meet mock connection expectations with query statement '%s'.", statement)
	}
}

func TestExecuteQueryHandlesEmptyStatementArg(t *testing.T) {
	mock := createMockConnection(t)

	defer mock.Close()

	rows, err := database.ExecuteQuery(mock, "")

	if err == nil {
		t.Fatal("Unable to catch error with empty 'statement' argument.")
	}

	if rows != nil {
		t.Fatalf("Actual rows response '%v' does not match expected nil rows response.", rows)
	}
}

func TestMapQueryResponseReturnsOne(t *testing.T) {
	expectedTitle := "Encanto"

	mock := createMockConnection(t)

	defer mock.Close()

	mock.ExpectQuery("SELECT \\* FROM movies WHERE id=1").
		WillReturnRows(pgxmock.
			NewRows([]string{"id", "title", "tagline", "description", "release_date", "runtime", "image", "reference"}).
			AddRow(1, "Encanto", "", "", "", 0, "", 812))

	statement, _ := database.CreateQuery("*", "movies", "id=1", "")
	rows, _ := database.ExecuteQuery(mock, statement)
	response, err := database.MapQueryResponse[model.MovieFragment](rows)

	if err != nil {
		t.Fatalf("Unable to map query response '%v' to internal movie fragment model: %v\n", response, err)
	}

	if len(response) != 1 {
		t.Fatal("Unable to map query response to exactly one movie fragment model.")
	}

	if response[0].Title != expectedTitle {
		t.Fatalf("Actual title '%s' does not match expected title '%s'.", response[0].Title, expectedTitle)
	}
}

func TestMapQueryResponseReturnsMany(t *testing.T) {
	mock := createMockConnection(t)

	defer mock.Close()

	mock.ExpectQuery("SELECT \\* FROM genres WHERE name='Action' OR name='Animation'").
		WillReturnRows(pgxmock.
			NewRows([]string{"id", "name", "reference"}).
			AddRow(1, "Action", 1).
			AddRow(2, "Animation", 2))

	statement, _ := database.CreateQuery("*", "genres", "name='Action' OR name='Animation'", "")
	rows, _ := database.ExecuteQuery(mock, statement)
	response, err := database.MapQueryResponse[model.MovieGenreFragment](rows)

	if err != nil {
		t.Fatalf("Unable to map query response '%v' to internal movie genre fragment model: %v\n", response, err)
	}

	if len(response) != 2 {
		t.Fatal("Unable to map query response to exactly two movie genre fragment models.")
	}

	if response[0].Name != "Action" || response[1].Name != "Animation" {
		t.Fatalf("Actual genre names '%s' and '%s' do not match expected names '%s' and '%s'.", response[0].Name, response[1].Name, "Action", "Animation")
	}
}

func TestMapQueryResponseReturnsNone(t *testing.T) {
	mock := createMockConnection(t)

	defer mock.Close()

	mock.ExpectQuery("SELECT \\* FROM production_companies WHERE id=4").
		WillReturnRows(pgxmock.
			NewRows([]string{"id", "name", "image", "reference"}))

	statement, _ := database.CreateQuery("*", "production_companies", "id=4", "")
	rows, _ := database.ExecuteQuery(mock, statement)
	response, err := database.MapQueryResponse[model.MovieProductionCompanyFragment](rows)

	if err != nil {
		t.Fatalf("Unable to map query response '%v' to internal movie production company fragment model: %v\n", response, err)
	}

	if len(response) != 0 {
		t.Fatal("Unable to map query response to exactly zero movie production company fragment models.")
	}
}

func TestMapQueryResponseHandlesTypeMismatch(t *testing.T) {
	mock := createMockConnection(t)

	defer mock.Close()

	mock.ExpectQuery("SELECT \\* FROM production_companies WHERE id=1").
		WillReturnRows(pgxmock.
			NewRows([]string{"id", "name", "image", "reference"}).
			AddRow(1, nil, nil, nil))

	statement, _ := database.CreateQuery("*", "production_companies", "id=1", "")
	rows, _ := database.ExecuteQuery(mock, statement)
	response, err := database.MapQueryResponse[model.MovieGenreFragment](rows)

	if err == nil {
		t.Fatal("Unable to catch error with type mismatch error due to 'nil' values in database.")
	}

	if len(response) != 0 {
		t.Fatalf("Actual response '%v' does not match expected empty response.", response)
	}
}

func createMockConnection(t *testing.T) pgxmock.PgxPoolIface {
	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Unable to create mock database pool connection: %v\n", err)
	}

	return mock
}

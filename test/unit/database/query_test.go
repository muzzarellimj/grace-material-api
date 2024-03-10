package database_test

import (
	"testing"

	"github.com/muzzarellimj/grace-material-api/database"
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

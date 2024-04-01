package service_test

import (
	"testing"

	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/movie"
	"github.com/pashagolub/pgxmock/v3"
)

func TestFetchFragmentSliceReturnsOne(t *testing.T) {
	mock := createMockConnection(t)

	defer mock.Close()

	mock.ExpectQuery("SELECT \\* FROM movies WHERE id=1").
		WillReturnRows(pgxmock.
			NewRows([]string{"id", "title", "tagline", "description", "release_date", "runtime", "image", "reference"}).
			AddRow(1, "", "", "", "", 0, "", 0))

	_, err := service.FetchFragmentSlice[model.MovieFragment](mock, database.TableMovieFragments, "id=1")

	if err != nil {
		t.Fatalf("Unable to fetch fragment slice: %v\n", err)
	}

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatalf("Mock connection expectations were not met: %v\n", err)
	}
}

func TestFetchFragmentSliceReturnsMany(t *testing.T) {
	mock := createMockConnection(t)

	defer mock.Close()

	mock.ExpectQuery("SELECT \\* FROM genres WHERE name='Action' OR name='Animation'").
		WillReturnRows(pgxmock.
			NewRows([]string{"id", "name", "reference"}).
			AddRow(1, "Action", 0).
			AddRow(2, "Animation", 0))

	_, err := service.FetchFragmentSlice[model.MovieGenreFragment](mock, database.TableMovieGenreFragments, "name='Action' OR name='Animation'")

	if err != nil {
		t.Fatalf("Unable to fetch fragment slice: %v\n", err)
	}

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatalf("Mock connection expectations were not met: %v\n", err)
	}
}

func TestFetchFragmentSliceReturnsNone(t *testing.T) {
	mock := createMockConnection(t)

	defer mock.Close()

	mock.ExpectQuery("SELECT \\* FROM production_companies WHERE id=4").
		WillReturnRows(pgxmock.
			NewRows([]string{"id", "name", "image", "reference"}))

	_, err := service.FetchFragmentSlice[model.MovieProductionCompanyFragment](mock, database.TableMovieProductionCompanyFragments, "id=4")

	if err != nil {
		t.Fatalf("Unable to fetch fragment slice: %v\n", err)
	}

	err = mock.ExpectationsWereMet()

	if err != nil {
		t.Fatalf("Mock connection expectations were not met: %v\n", err)
	}
}

func createMockConnection(t *testing.T) pgxmock.PgxPoolIface {
	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("Unable to create mock database pool connection: %v\n", err)
	}

	return mock
}

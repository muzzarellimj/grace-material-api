package api_test

import (
	"testing"

	"github.com/joho/godotenv"
	api "github.com/muzzarellimj/grace-material-api/internal/api/third_party/igdb.com"
	model "github.com/muzzarellimj/grace-material-api/internal/model/third_party/igdb.com"
)

func TestIGDBGetResourceReturnsGame(t *testing.T) {
	godotenv.Load("../../../../../../.env")

	expectedTitle := "The Last of Us"

	actual, err := api.IGDBGetResource[model.IGDBGameResponse](api.IGDBEndpointGame, "fields id,cover.*,first_release_date,franchises.*,genres.*,involved_companies.*,name,platforms.*,storyline,summary; where id = 1009;")

	if err != nil {
		t.Fatalf("Unable to execute request to get IGDB resource: %v\n", err)
	}

	if actual.Title != expectedTitle {
		t.Fatalf("Actual title '%s' does not match expected title '%s'.", actual.Title, expectedTitle)
	}
}

func TestIGDBGetResourceReturnsEmpty(t *testing.T) {
	godotenv.Load("../../../../../../.env")

	actual, err := api.IGDBGetResource[model.IGDBGameResponse](api.IGDBEndpointGame, "fields id,cover.*,first_release_date,franchises.*,genres.*,involved_companies.*,name,platforms.*,storyline,summary; where id = -1;")

	if err != nil {
		t.Fatalf("Unable to execute request to get IGDB resource: %v\n", err)
	}

	if actual.ID != 0 {
		t.Fatalf("Actual numeric identifier '%d' does not match expected zero numeric identifier", actual.ID)
	}
}

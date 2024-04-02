package api_test

import (
	"fmt"
	"testing"

	api "github.com/muzzarellimj/grace-material-api/internal/api/third_party/openlibrary.org"
)

func TestOLGetAuthorReturnsOLAuthorResponse(t *testing.T) {
	id := "OL368638A"

	expectedId := fmt.Sprint("/authors/", id)
	expectedName := "Andrzej Sapkowski"

	author, err := api.OLGetAuthor(id)

	if err != nil {
		t.Fatalf("Unable to retrieve and decode author '%s': %v\n", id, err)
	}

	if author.ID != expectedId {
		t.Fatalf("Actual author identifier '%s' does not match expected author identifier '%s'.", author.ID, expectedId)
	}

	if author.Name != expectedName {
		t.Fatalf("Actual author name '%s' does not match expected author name '%s'.", author.Name, expectedName)
	}
}

func TestOLGetAuthorReturnsEmptyOLAuthorResponse(t *testing.T) {
	id := "OL0A"

	author, err := api.OLGetAuthor(id)

	if err != nil {
		t.Fatalf("Unable to retrieve and decode author '%s': %v\n", id, err)
	}

	if author.ID != "" {
		t.Fatalf("Actual author identifier '%s' does not match expected empty author identifier.", author.ID)
	}
}

func TestOLGetAuthorReturnsError(t *testing.T) {
	id := ""

	author, err := api.OLGetAuthor(id)

	if err == nil {
		t.Fatal("Unable to catch error with missing 'id' argument.")
	}

	if author.ID != "" {
		t.Fatalf("Actual author identifier '%s' does not match expected empty author identifier.", author.ID)
	}
}

func TestOLGetEditionReturnsOLEditionResponse(t *testing.T) {
	id := "9780316452465"

	expectedId := "/books/OL37765857M"
	expectedTitle := "The Last Wish"

	edition, err := api.OLGetEdition(id)

	if err != nil {
		t.Fatalf("Unable to retrieve and decode edition '%s': %v\n", id, err)
	}

	if edition.ID != expectedId {
		t.Fatalf("Actual edition identifier '%s' does not match expected edition identifier '%s'.", edition.ID, expectedId)
	}

	if edition.Title != expectedTitle {
		t.Fatalf("Actual edition title '%s' does not match expected edition title '%s'.", edition.Title, expectedTitle)
	}
}

func TestOLGetEditionReturnsEmptyOLEditionResponse(t *testing.T) {
	id := "OL0M"

	edition, err := api.OLGetEdition(id)

	if err != nil {
		t.Fatalf("Unable to retrieve and decode edition '%s': %v\n", id, err)
	}

	if edition.ID != "" {
		t.Fatalf("Actual edition identifier '%s' does not match expected empty edition identifier.", edition.ID)
	}
}

func TestOLGetEditionReturnsError(t *testing.T) {
	id := ""

	edition, err := api.OLGetEdition(id)

	if err == nil {
		t.Fatal("Unable to catch error with missing 'id' argument.")
	}

	if edition.ID != "" {
		t.Fatalf("Actual edition identifier '%s' does not match expected empty edition identifier.", edition.ID)
	}
}

func TestOLGetWorkReturnsOLWorkResponse(t *testing.T) {
	id := "OL27691456W"

	expectedId := fmt.Sprint("/works/", id)
	expectedTitle := "The Last Wish"

	work, err := api.OLGetWork(id)

	if err != nil {
		t.Fatalf("Unable to retrieve and decode work '%s': %v\n", id, err)
	}

	if work.ID != expectedId {
		t.Fatalf("Actual work identifier '%s' does not match expected work identifier '%s'.", work.ID, expectedId)
	}

	if work.Title != expectedTitle {
		t.Fatalf("Actual work title '%s' does not match expected work title '%s'.", work.Title, expectedTitle)
	}
}

func TestOLGetWorkReturnsEmptyOLWorkResponse(t *testing.T) {
	id := "OL0W"

	work, err := api.OLGetWork(id)

	if err != nil {
		t.Fatalf("Unable to retrieve and decode work '%s': %v\n", id, err)
	}

	if work.ID != "" {
		t.Fatalf("Actual work identifier '%s' does not match expected empty work identifier.", work.ID)
	}
}

func TestOLGetWorkReturnsError(t *testing.T) {
	id := ""

	work, err := api.OLGetWork(id)

	if err == nil {
		t.Fatal("Unable to catch error with missing 'id' argument.")
	}

	if work.ID != "" {
		t.Fatalf("Actual work identifier '%s' does not match expected empty work identifier.", work.ID)
	}
}

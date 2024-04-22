package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	model "github.com/muzzarellimj/grace-material-api/internal/model/third_party/openlibrary.org"
	"github.com/muzzarellimj/grace-material-api/internal/util"
)

const (
	OLBase            = "https://openlibrary.org"
	OLEndpointAuthor  = "/authors"
	OLEndpointEdition = "/isbn"
	OLEndpointSearch  = "/search"
	OLEndpointWork    = "/works"
)

// Get an authority control figure which can include a name, viographical information, and image, among other items.
//
// Return: decoded author response and nil with sucess, empty author response and error without.
func OLGetAuthor(id string) (model.OLAuthorResponse, error) {
	var zero model.OLAuthorResponse

	if id == "" {
		err := errors.New("unable to process request with missing 'id' arg")

		fmt.Fprintf(os.Stderr, "Unable to process author request due to missing identifier argument.")

		return zero, err
	}

	path, err := util.CreateRequestPath(OLBase, OLEndpointAuthor, fmt.Sprint(id, ".json"), map[string]string{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request path to '%s%s': %v\n", OLBase, OLEndpointAuthor, err)

		return zero, err
	}

	request, err := util.CreateRequest(http.MethodGet, path, []byte{}, map[string]string{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create '%s' request to '%s': %v\n", http.MethodGet, path, err)

		return zero, err
	}

	response, err := util.ExecuteRequest(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute '%s' request to '%s': %v\n", request.Method, request.URL.String(), err)

		return zero, err
	}

	if response.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stdout, "Unable to find author matching identifier '%s'.", id)

		return zero, nil
	}

	var author model.OLAuthorResponse

	err = json.NewDecoder(response.Body).Decode(&author)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode response as author model: %v\n", err)

		return zero, err
	}

	return author, nil
}

// Get the material record unique to an ISBN where data may differ between other editions of the same work.
//
// Return: decoded edition response and nil with success, empty edition response and error without.
func OLGetEdition(id string) (model.OLEditionResponse, error) {
	var zero model.OLEditionResponse

	if id == "" {
		err := errors.New("unable to process request with missing 'id' arg")

		fmt.Fprintf(os.Stderr, "Unable to process edition request due to missing identifier argument.")

		return zero, err
	}

	path, err := util.CreateRequestPath(OLBase, OLEndpointEdition, fmt.Sprint(id, ".json"), map[string]string{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request path to '%s%s': %v\n", OLBase, OLEndpointEdition, err)

		return zero, err
	}

	request, err := util.CreateRequest(http.MethodGet, path, []byte{}, map[string]string{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create '%s' request to '%s': %v\n", http.MethodGet, path, err)

		return zero, err
	}

	response, err := util.ExecuteRequest(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute '%s' request to '%s': %v\n", request.Method, request.URL.String(), err)

		return zero, err
	}

	if response.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stdout, "Unable to find edition matching identifier '%s'.", id)

		return zero, nil
	}

	var edition model.OLEditionResponse

	err = json.NewDecoder(response.Body).Decode(&edition)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode response as edition model: %v\n", err)

		return zero, err
	}

	return edition, nil
}

// Get the logical, more abstract collection record which can include data shared across all independent editions that
// may differ in language, version, cover, etc.
//
// Return: decoded work response and nil with success, empty work response and error without.
func OLGetWork(id string) (model.OLWorkResponse, error) {
	var zero model.OLWorkResponse

	if id == "" {
		err := errors.New("unable to process request with missing 'id' arg")

		fmt.Fprintf(os.Stderr, "Unable to process work request due to missing identifier argument.")

		return zero, err
	}

	path, err := util.CreateRequestPath(OLBase, OLEndpointWork, fmt.Sprint(id, ".json"), map[string]string{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request path to '%s%s': %v\n", OLBase, OLEndpointWork, err)

		return zero, err
	}

	request, err := util.CreateRequest(http.MethodGet, path, []byte{}, map[string]string{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create '%s' request to '%s': %v\n", http.MethodGet, path, err)

		return zero, err
	}

	response, err := util.ExecuteRequest(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute '%s' request to '%s': %v\n", request.Method, request.URL.String(), err)

		return zero, err
	}

	if response.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stdout, "Unable to find work matching identifier '%s'.", id)

		return zero, nil
	}

	var work model.OLWorkResponse

	err = json.NewDecoder(response.Body).Decode(&work)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode response as work model: %v\n", err)

		return zero, err
	}

	return work, nil
}

func OLSearchBook(query string) (model.OLBookSearchResponse, error) {
	var zero model.OLBookSearchResponse

	if query == "" {
		err := errors.New("unable to process request with missing 'query' arg")

		fmt.Fprintf(os.Stderr, "Unable to process book search request due to missing query argument.\n")

		return zero, err
	}

	path, err := util.CreateRequestPath(OLBase, fmt.Sprint(OLEndpointSearch, ".json"), "", map[string]string{"q": fmt.Sprint(query, " language:eng")})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request path to '%s%s': %v\n", OLBase, OLEndpointSearch, err)

		return zero, err
	}

	request, err := util.CreateRequest(http.MethodGet, path, []byte{}, map[string]string{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create '%s' request to '%s': %v\n", http.MethodGet, path, err)

		return zero, err
	}

	response, err := util.ExecuteRequest(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute '%s' request to '%s': %v\n", request.Method, request.URL.String(), err)

		return zero, err
	}

	if response.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stdout, "Unable to find results matching query '%s'.\n", query)

		return zero, nil
	}

	var model model.OLBookSearchResponse

	err = json.NewDecoder(response.Body).Decode(&model)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode response as search response model: %v\n", err)

		return zero, err
	}

	return model, nil
}

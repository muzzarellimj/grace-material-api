package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/muzzarellimj/grace-material-api/pkg/util"
)

const (
	IGDBEndpointCompany         = "/v4/companies"
	IGDBEndpointCover           = "/v4/covers"
	IGDBEndpointFranchise       = "/v4/franchises"
	IGDBEndpointGame            = "/v4/games"
	IGDBEndpointGenre           = "/v4/genres"
	IGDBEndpointInvolvedCompany = "/v4/involved_companies"
	IGDBEndpointPlatform        = "/v4/platforms"
)

// Get an IGDB resource slice with a provided model to decode to and an Apicalypse-compliant constraint.
//
// Return: decoded model slice and nil with success, empty model slice and error without.
func IGDBGetResourceSlice[M interface{}](endpoint string, constraint string) ([]M, error) {
	var zero []M

	path, err := util.CreateRequestPath(os.Getenv("AWS_PROXY_HOST"), endpoint, "", map[string]string{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request path to AWS IGDB proxy at endpoint '%s': %v\n", endpoint, err)

		return zero, err
	}

	request, err := util.CreateRequest(http.MethodPost, path, []byte(constraint), map[string]string{
		"x-api-key": os.Getenv("AWS_PROXY_API_KEY"),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request to AWS IGDB proxy at endpoint '%s': %v\n", endpoint, err)

		return zero, err
	}

	response, err := util.ExecuteRequest(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute request to AWS IGDB proxy at endpoint '%s': %v\n", endpoint, err)

		return zero, err
	}

	var resourceSlice []M

	err = json.NewDecoder(response.Body).Decode(&resourceSlice)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode response to provided resource model: %v\n", err)

		return zero, err
	}

	return resourceSlice, nil
}

// Get one IGDB resource with a provided model to decode to and an Apicalyps-compliant constraint.
//
// Return: decoded model and nil with success, empty model and error without.
func IGDBGetResource[M interface{}](endpoint string, constraint string) (M, error) {
	var zero M

	resourceSlice, err := IGDBGetResourceSlice[M](endpoint, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get and decode initial resource slice: %v\n", err)

		return zero, err
	}

	if len(resourceSlice) > 0 {
		return resourceSlice[0], nil
	}

	return zero, nil
}

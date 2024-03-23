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

// Get an IGDB resource with a provided model and Apicalypse-compliant constraint.
//
// Return: decoded game model and nil with success, empty game response and error without.
func IGDBGetResource[M interface{}](endpoint string, constraint string) (M, error) {
	var zero M

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

	if len(resourceSlice) > 0 {
		return resourceSlice[0], nil
	}

	return zero, nil
}

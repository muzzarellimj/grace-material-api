package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	model "github.com/muzzarellimj/grace-material-api/internal/model/third_party/twitch.tv"
	"github.com/muzzarellimj/grace-material-api/internal/util"
)

const (
	base                   = "https://id.twitch.tv"
	endpointAuthentication = "/oauth2/token"
)

func TTVAuthenticateRequest() (string, error) {
	path, err := util.CreateRequestPath(base, endpointAuthentication, "", map[string]string{
		"client_id":     os.Getenv("IGDB_CLIENT_ID"),
		"client_secret": os.Getenv("IGDB_CLIENT_SECRET"),
		"grant_type":    "client_credentials",
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request path to '%s%s': %v\n", base, endpointAuthentication, err)

		return "", err
	}

	request, err := util.CreateRequest(http.MethodPost, path, []byte{}, map[string]string{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create '%s' request to '%s': %v\n", http.MethodPost, path, err)

		return "", err
	}

	response, err := util.ExecuteRequest(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute '%s' request to '%s': %v\n", request.Method, request.URL.String(), err)

		return "", err
	}

	var authentication model.TTVAuthenticationResponse

	err = json.NewDecoder(response.Body).Decode(&authentication)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode response as authentication response model: %v\n", err)

		return "", err
	}

	if authentication.AccessToken == "" || authentication.TokenType != "bearer" {
		fmt.Fprintf(os.Stderr, "Unable to parse expected properties from authentication response model: %v\n", err)

		return "", err
	}

	return authentication.AccessToken, nil
}

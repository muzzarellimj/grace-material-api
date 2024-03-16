package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	model "github.com/muzzarellimj/grace-material-api/pkg/model/third_party/themoviedb.org"
	"github.com/muzzarellimj/grace-material-api/pkg/util"
)

const (
	TMDBBase          = "https://api.themoviedb.org/3"
	TMDBEndpointMovie = "/movie/"
)

// Get the top-level details of a movie with a provided numeric identifier.
//
// Return: decoded movie detail response and nil with success, empty movie detail response and error without.
func TMDBGetMovie(id int) (model.TMDBMovieDetailResponse, error) {
	path, err := util.CreateRequestPath(TMDBBase, TMDBEndpointMovie, fmt.Sprint(id), map[string]string{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request path to '%s%s': %v\n", TMDBBase, TMDBEndpointMovie, err)

		return model.TMDBMovieDetailResponse{}, err
	}

	request, err := util.CreateRequest(http.MethodGet, path, map[string]string{"Authorization": fmt.Sprint("Bearer ", os.Getenv("TMDB_API_KEY"))})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create '%s' request to '%s': %v\n", http.MethodGet, path, err)

		return model.TMDBMovieDetailResponse{}, err
	}

	response, err := util.ExecuteRequest(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute '%s' request to '%s': %v\n", request.Method, request.URL.String(), err)

		return model.TMDBMovieDetailResponse{}, err
	}

	var movie model.TMDBMovieDetailResponse

	err = json.NewDecoder(response.Body).Decode(&movie)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode response as movie detail model: %v\n", err)

		return model.TMDBMovieDetailResponse{}, err
	}

	return movie, nil
}

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/internal/api/movie/helper"
	tapi "github.com/muzzarellimj/grace-material-api/internal/api/third_party/themoviedb.org"
)

func HandleGetMovie(context *gin.Context) {
	id, err := strconv.Atoi(context.Query("id"))

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	movie, message, err := fetchMovie(fmt.Sprintf("id=%d", id))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": message,
		})

		return
	}

	if movie.ID != 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   movie,
		})

		return
	}

	context.Status(http.StatusNoContent)
}

func HandlePostMovie(context *gin.Context) {
	id, err := strconv.Atoi(context.Query("id"))

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	existingMovie, message, err := fetchMovie(fmt.Sprintf("reference=%d", id))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": message,
		})

		return
	}

	if existingMovie.ID != 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   existingMovie,
		})

		return
	}

	tmdbMovie, err := tapi.TMDBGetMovie(id)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": message,
		})

		return
	}

	if tmdbMovie.ID == 0 {
		context.Status(http.StatusNoContent)

		return
	}

	insertedMovieId, err := storeMovie(tmdbMovie)

	if err != nil || insertedMovieId == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": message,
		})

		return
	}

	context.IndentedJSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": fmt.Sprintf("Movie inserted with numeric identifier '%d'.", insertedMovieId),
	})
}

func HandleGetMovieSearch(context *gin.Context) {
	query := context.Query("query")

	if query == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid search term '%s' provided in query parameter 'query'.", context.Query("query")),
		})

		return
	}

	results, err := tapi.TMDBSearchMovie(query)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Unable to fetch movie metadata and map to supported data structure.",
		})

		return
	}

	if len(results.Results) > 0 {
		mappedResults := helper.MapSearchResultSlice(results.Results)

		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   mappedResults,
		})

		return
	}

	context.Status(http.StatusNoContent)
}

package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/internal/api/movie/helper"
	tapi "github.com/muzzarellimj/grace-material-api/internal/api/third_party/themoviedb.org"
)

const errorMessage string = "Unable to fetch movie metadata and map to supported data structure."

func HandleGetMovie(context *gin.Context) {
	idArg := context.Query("id")

	if len(idArg) == 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier argument '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	idSlice := strings.Split(idArg, ",")

	if len(idSlice) == 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier argument '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	var constraintSlice []string

	for _, id := range idSlice {
		constraintSlice = append(constraintSlice, fmt.Sprintf("id=%s", id))
	}

	movieSlice, errorSlice := helper.FetchMovieSlice(constraintSlice)

	if len(errorSlice) != 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": errorMessage,
		})

		return
	}

	if len(movieSlice) == 0 {
		context.Status(http.StatusNoContent)

		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   movieSlice,
	})
}

func HandlePostMovie(context *gin.Context) {
	idArg := context.Query("id")

	if len(idArg) == 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier argument '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	existingMovie, err := helper.FetchMovie(fmt.Sprintf("reference=%s", idArg))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	if existingMovie.ID != 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data": map[string]any{
				"id": existingMovie.ID,
			},
		})

		return
	}

	movie, err := tapi.TMDBGetMovie(idArg)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	if movie.ID == 0 {
		context.Status(http.StatusNoContent)

		return
	}

	storedMovieId, err := helper.ProcessMovieStorage(movie)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	context.IndentedJSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"data": map[string]any{
			"id": storedMovieId,
		},
	})
}

// func OldHandlePostMovie(context *gin.Context) {
// 	id, err := strconv.Atoi(context.Query("id"))

// 	if err != nil {
// 		context.IndentedJSON(http.StatusBadRequest, gin.H{
// 			"status":  http.StatusBadRequest,
// 			"message": fmt.Sprintf("Invalid material identifier '%s' provided in query parameter 'id'.", context.Query("id")),
// 		})

// 		return
// 	}

// 	existingMovie, err := helper.FetchMovie(fmt.Sprintf("reference=%d", id))

// 	if err != nil {
// 		context.IndentedJSON(http.StatusInternalServerError, gin.H{
// 			"status":  http.StatusInternalServerError,
// 			"message": errorMessage,
// 		})

// 		return
// 	}

// 	if existingMovie.ID != 0 {
// 		context.IndentedJSON(http.StatusOK, gin.H{
// 			"status": http.StatusOK,
// 			"data":   existingMovie,
// 		})

// 		return
// 	}

// 	tmdbMovie, err := tapi.TMDBGetMovie(id)

// 	if err != nil {
// 		context.IndentedJSON(http.StatusInternalServerError, gin.H{
// 			"status":  http.StatusInternalServerError,
// 			"message": errorMessage,
// 		})

// 		return
// 	}

// 	if tmdbMovie.ID == 0 {
// 		context.Status(http.StatusNoContent)

// 		return
// 	}

// 	insertedMovieId, err := storeMovie(tmdbMovie)

// 	if err != nil || insertedMovieId == 0 {
// 		context.IndentedJSON(http.StatusInternalServerError, gin.H{
// 			"status":  http.StatusInternalServerError,
// 			"message": errorMessage,
// 		})

// 		return
// 	}

// 	context.IndentedJSON(http.StatusCreated, gin.H{
// 		"status":  http.StatusCreated,
// 		"message": fmt.Sprintf("Movie inserted with numeric identifier '%d'.", insertedMovieId),
// 	})
// }

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

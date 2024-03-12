package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/database"
	"github.com/muzzarellimj/grace-material-api/model"
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

	movie, message, err := fetchMovie("id", id)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": message,
		})

		return
	}

	if movie != (model.Movie{}) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   movie,
		})

		return
	}

	context.Status(http.StatusNoContent)
}

func fetchMovie(field string, id int) (model.Movie, string, error) {
	statement, err := database.CreateQuery(
		"m.id, m.title, m.description, m.tagline, ARRAY_AGG(DISTINCT g.name) as genres, ARRAY_AGG(DISTINCT p.name) as production_companies, m.release_date, m.runtime, m.image, m.reference_imdb, m.reference_tmdb",
		"movies m",
		fmt.Sprintf("m.%s=%v", field, id),
		"m.id",
		"JOIN movies_genres mg ON m.id = mg.movie",
		"JOIN genres g ON g.id = mg.genre",
		"JOIN movies_production_companies mp ON m.id = mp.movie",
		"JOIN production_companies p ON p.id = mp.production_company",
	)

	if err != nil {
		message := "Unable to create database query statement."

		return model.Movie{}, message, err
	}

	rows, err := database.ExecuteQuery(database.MovieConnection, statement)

	if err != nil {
		message := "Unable to execute database query."

		return model.Movie{}, message, err
	}

	response, err := database.MapResponse[model.Movie](rows)

	if err != nil {
		message := "Unable to map database response to supported data structure."

		return model.Movie{}, message, err
	}

	if len(response) == 0 {
		return model.Movie{}, "", nil
	}

	return response[0], "", nil
}

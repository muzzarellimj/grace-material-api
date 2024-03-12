package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/database"
	"github.com/muzzarellimj/grace-material-api/model"
)

func HandleGetMovie(context *gin.Context) {
	id := context.Query("id")

	if id == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid numeric material identifier provided in query parameters.",
		})

		return
	}

	statement, err := database.CreateQuery(
		"m.id, m.title, m.description, m.tagline, ARRAY_AGG(DISTINCT g.name) as genres, ARRAY_AGG(DISTINCT p.name) as production_companies, m.release_date, m.runtime, m.image, m.reference_imdb, m.reference_tmdb",
		"movies m",
		fmt.Sprint("m.id=", id),
		"m.id",
		"JOIN movies_genres mg ON m.id = mg.movie",
		"JOIN genres g ON g.id = mg.genre",
		"JOIN movies_production_companies mp ON m.id = mp.movie",
		"JOIN production_companies p ON p.id = mp.production_company",
	)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Unable to create PostgreSQL query statement.",
		})

		return
	}

	rows, err := database.ExecuteQuery(database.MovieConnection, statement)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Unable to execute PostgreSQL query.",
		})

		return
	}

	response, err := database.MapResponse[model.Movie](rows)

	if err != nil || len(response) > 1 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Unable to map response to a supported data structure.",
		})

		return
	}

	if len(response) > 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   response[0],
		})

		return
	}

	context.IndentedJSON(http.StatusNoContent, gin.H{
		"status":  http.StatusNoContent,
		"message": "No material with the provided numeric material identifier could be found.",
	})
}

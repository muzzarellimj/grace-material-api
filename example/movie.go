package example

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Overview string `json:"overview"`
}

var movies = []Movie{
	{ID: 568124, Title: "Encanto", Overview: "The tale of an extraordinary family, the Madrigals, who live hidden in the mountains of Colombia, in a magical house, in a vibrant town, in a wondrous, charmed place called an Encanto. The magic of the Encanto has blessed every child in the family—every child except one, Mirabel. But when she discovers that the magic surrounding the Encanto is in danger, Mirabel decides that she, the only ordinary Madrigal, might just be her exceptional family's last hope."},
	{ID: 38757, Title: "Tangled", Overview: "When the kingdom's most wanted-and most charming-bandit Flynn Rider hides out in a mysterious tower, he's taken hostage by Rapunzel, a beautiful and feisty tower-bound teen with 70 feet of magical, golden hair. Flynn's curious captor, who's looking for her ticket out of the tower where she's been locked away for years, strikes a deal with the handsome thief and the unlikely duo sets off on an action-packed escapade, complete with a super-cop horse, an over-protective chameleon and a gruff gang of pub thugs."},
}

func FetchMovies(context *gin.Context) {
	limit, error := strconv.Atoi(context.Query("limit"))

	if error != nil {
		limit = -1
	}

	if len(movies) == 0 {
		context.IndentedJSON(http.StatusNoContent, nil)

		return
	}

	if limit > 0 {
		context.IndentedJSON(http.StatusOK, movies[0:limit])
	} else {
		context.IndentedJSON(http.StatusOK, movies)
	}
}

func FetchMovie(context *gin.Context) {
	id, error := strconv.Atoi(context.Param("id"))

	if error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Please provide a valid material identifier in the request path.",
		})

		return
	}

	for _, movie := range movies {
		if movie.ID == id {
			context.IndentedJSON(http.StatusOK, movie)

			return
		}
	}

	context.IndentedJSON(http.StatusNoContent, nil)
}

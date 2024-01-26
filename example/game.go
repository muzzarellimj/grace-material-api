package example

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Game struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Platforms []int  `json:"platforms"`
}

var games = []Game{
	{ID: 1942, Title: "The Witcher 3: Wild Hunt", Platforms: []int{6, 48, 49, 130, 167, 169}},
	{ID: 1905, Title: "Fortnite", Platforms: []int{6, 14, 34, 39, 48, 49, 130, 167, 169}},
}

func FetchGames(context *gin.Context) {
	limit, error := strconv.Atoi(context.Query("limit"))

	if error != nil {
		limit = -1
	}

	if len(games) == 0 {
		context.IndentedJSON(http.StatusNoContent, nil)

		return
	}

	if limit > 0 {
		context.IndentedJSON(http.StatusOK, games[0:limit])
	} else {
		context.IndentedJSON(http.StatusOK, games)
	}
}

func FetchGame(context *gin.Context) {
	id, error := strconv.Atoi(context.Param("id"))

	if error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Please provide a valid material identifier in the request path.",
		})

		return
	}

	for _, game := range games {
		if game.ID == id {
			context.IndentedJSON(http.StatusOK, game)

			return
		}
	}

	context.IndentedJSON(http.StatusNoContent, nil)
}

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/pkg/api/game/helper"
	IGDBAPI "github.com/muzzarellimj/grace-material-api/pkg/api/third_party/igdb.com"
	IGDBModel "github.com/muzzarellimj/grace-material-api/pkg/model/third_party/igdb.com"
)

const errorResponseMessage string = "Unable to fetch game metadata and map to supported data structure."

func HandleGetGame(context *gin.Context) {
	id, err := strconv.Atoi(context.Query("id"))

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	game, err := helper.FetchGame(fmt.Sprintf("id=%d", id))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	if game.ID != 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   game,
		})

		return
	}

	context.Status(http.StatusNoContent)
}

func HandlePostGame(context *gin.Context) {
	id, err := strconv.Atoi(context.Query("id"))

	if err != nil || id <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	existingGame, err := helper.FetchGame(fmt.Sprintf("reference=%d", id))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	if existingGame.ID != 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   existingGame,
		})

		return
	}

	game, err := IGDBAPI.IGDBGetResource[IGDBModel.IGDBGameResponse](IGDBAPI.IGDBEndpointGame, fmt.Sprintf("fields id,cover.*,first_release_date,franchises.*,genres.*,involved_companies.*,name,platforms.*,storyline,summary; where id=%d;", id))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	if game.ID == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	storedGameId, err := helper.ProcessGameStorage(game)

	if err != nil || storedGameId == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	context.IndentedJSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": fmt.Sprintf("Game stored with numeric identifier '%d'.", storedGameId),
	})
}

func HandleGetGameSearch(context *gin.Context) {
	query := context.Query("query")

	if query == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid search term '%s' provided in query parameter 'query'.", context.Query("query")),
		})

		return
	}

	results, err := IGDBAPI.IGDBGetResourceSlice[IGDBModel.IGDBGameSearchResponse](IGDBAPI.IGDBEndpointGame, fmt.Sprintf(`fields id,name,cover.*,first_release_date; search "%s"; where (status=0 | status=null) & category=0;`, query))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	if len(results) > 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   results,
		})

		return
	}

	context.Status(http.StatusNoContent)
}

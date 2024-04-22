package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/internal/api/game/helper"
	IGDBAPI "github.com/muzzarellimj/grace-material-api/internal/api/third_party/igdb.com"
	IGDBModel "github.com/muzzarellimj/grace-material-api/internal/model/third_party/igdb.com"
)

const errorMessage string = "Unable to fetch game metadata and map to supported data structure."

func HandleGetGame(context *gin.Context) {
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

	gameSlice, errSlice := helper.FetchGameSlice(constraintSlice)

	if len(errSlice) != 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	if len(gameSlice) == 0 {
		context.Status(http.StatusNoContent)

		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   gameSlice,
	})
}

func HandlePostGame(context *gin.Context) {
	idArg := context.Query("id")

	if len(idArg) == 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier argument '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	existingGame, err := helper.FetchGame(fmt.Sprintf("reference=%s", idArg))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	if existingGame.ID != 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data": map[string]any{
				"id": existingGame.ID,
			},
		})

		return
	}

	game, err := IGDBAPI.IGDBGetResource[IGDBModel.IGDBGameResponse](IGDBAPI.IGDBEndpointGame, fmt.Sprintf("fields id,cover.*,first_release_date,franchises.*,genres.*,involved_companies.*,name,platforms.*,storyline,summary; where id=%s;", idArg))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	if game.ID == 0 {
		context.Status(http.StatusNoContent)

		return
	}

	storedGameId, err := helper.ProcessGameStorage(game)

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
			"id": storedGameId,
		},
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
			"message": errorMessage,
		})

		return
	}

	if len(results) > 0 {
		mappedResults := helper.MapSearchResultSlice(results)

		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   mappedResults,
		})

		return
	}

	context.Status(http.StatusNoContent)
}

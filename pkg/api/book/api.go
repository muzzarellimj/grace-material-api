package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleGetBook(context *gin.Context) {
	id, err := strconv.Atoi(context.Query("id"))

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	book, err := fetchBook(fmt.Sprintf("id=%d", id))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Unable to fetch book metadata and map to a supported data structure.",
		})

		return
	}

	if book.ID != 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   book,
		})

		return
	}

	context.Status(http.StatusNoContent)
}

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/pkg/api/book/helper"
	api "github.com/muzzarellimj/grace-material-api/pkg/api/third_party/openlibrary.org"
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

func HandlePostBook(context *gin.Context) {
	errorResponseMessage := "Unable to fetch book metadata and map to supported data structure."

	id := helper.FormatISBN(context.Query("id"))

	if id == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	existingBook, err := fetchBook(fmt.Sprintf("isbn13='%s' OR edition_reference='%s'", id, id))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	if existingBook.ID != 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   existingBook,
		})

		return
	}

	edition, err := api.OLGetEdition(id)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	if len(edition.Works) == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	work, err := api.OLGetWork(helper.ExtractResourceId(edition.Works[0].ID))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	insertedBookId, err := storeBook(edition, work)

	if err != nil || insertedBookId == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	context.IndentedJSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": fmt.Sprintf("Book inserted with numeric identifier '%d'.", insertedBookId),
	})
}

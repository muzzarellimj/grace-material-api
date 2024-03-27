package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/pkg/api/book/helper"
	api "github.com/muzzarellimj/grace-material-api/pkg/api/third_party/openlibrary.org"
	model "github.com/muzzarellimj/grace-material-api/pkg/model/book"
)

const errorResponseMessage string = "Unable to fetch book metadata and map to supported data structure."

func HandleGetBook(context *gin.Context) {
	id, err := strconv.Atoi(context.Query("id"))

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	book, err := helper.FetchBook(fmt.Sprintf("id=%d", id))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
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
	id := helper.FormatISBN(context.Query("id"))

	if id == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	existingBook, err := helper.FetchBook(fmt.Sprintf("isbn13='%s' OR edition_reference='%s'", id, id))

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

	storedBookId, err := helper.ProcessBookStorage(edition, work)

	if err != nil || storedBookId == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	context.IndentedJSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": fmt.Sprintf("Book stored with numeric identifier '%d'.", storedBookId),
	})
}

func HandleGetBookSearch(context *gin.Context) {
	query := helper.FormatISBN(context.Query("query"))

	if query == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid value '%s' provided in query parameter 'query'.", context.Query("query")),
		})

		return
	}

	edition, err := api.OLGetEdition(query)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorResponseMessage,
		})

		return
	}

	if edition.ID == "" {
		context.Status(http.StatusNoContent)

		return
	}

	result := model.BookSearchResult{
		ID:          helper.ExtractResourceId(edition.ID),
		Title:       edition.Title,
		PublishDate: edition.PublishDate,
		ISBN10:      helper.ExtractISBN(edition.ISBN10),
		ISBN13:      helper.ExtractISBN(edition.ISBN13),
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   result,
	})
}

package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/internal/api/book/helper"
	api "github.com/muzzarellimj/grace-material-api/internal/api/third_party/openlibrary.org"
	model "github.com/muzzarellimj/grace-material-api/internal/model/book"
	"github.com/muzzarellimj/grace-material-api/internal/util"
)

const errorMessage string = "Unable to fetch book metadata and map to supported data structure."

func HandleGetBook(context *gin.Context) {
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

	bookSlice, errSlice := helper.FetchBookSlice(constraintSlice)

	if len(errSlice) != 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	if len(bookSlice) == 0 {
		context.Status(http.StatusNoContent)

		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   bookSlice,
	})
}

func HandlePostBook(context *gin.Context) {
	idArg := context.Query("id")

	if len(idArg) == 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid material identifier argument '%s' provided in query parameter 'id'.", context.Query("id")),
		})

		return
	}

	idArg = helper.FormatISBN(context.Query("id"))

	existingBook, err := helper.FetchBook(fmt.Sprintf("isbn13='%s' OR edition_reference='%s'", idArg, idArg))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	if existingBook.ID != 0 {
		context.IndentedJSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data": map[string]any{
				"id": existingBook.ID,
			},
		})

		return
	}

	edition, err := api.OLGetEdition(idArg)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	if len(edition.Works) == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	work, err := api.OLGetWork(helper.ExtractResourceId(edition.Works[0].ID))

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	storedBookId, err := helper.ProcessBookStorage(edition, work)

	if err != nil || storedBookId == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	context.IndentedJSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"data": map[string]any{
			"id": storedBookId,
		},
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
			"message": errorMessage,
		})

		return
	}

	if edition.ID == "" {
		context.Status(http.StatusNoContent)

		return
	}

	results := []model.BookSearchResult{
		{
			ID:          helper.ExtractResourceId(edition.ID),
			Title:       edition.Title,
			PublishDate: util.ParseDateTime(edition.PublishDate),
			ISBN10:      helper.ExtractISBN(edition.ISBN10),
			ISBN13:      helper.ExtractISBN(edition.ISBN13),
		},
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   results,
	})
}

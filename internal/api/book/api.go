package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muzzarellimj/grace-material-api/internal/api/book/helper"
	OLAPI "github.com/muzzarellimj/grace-material-api/internal/api/third_party/openlibrary.org"
	model "github.com/muzzarellimj/grace-material-api/internal/model/book"
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

func HandlePutBook(context *gin.Context) {
	var book model.BookFragment

	err := context.BindJSON(&book)

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Unable to bind request JSON body to book model.",
		})

		return
	}

	id, err := helper.UpdateBookFragment(book)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Unable to update book fragment.",
		})

		return
	}

	if id == 0 {
		context.Status(http.StatusNoContent)

		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": map[string]any{
			"id": id,
		},
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

	edition, err := OLAPI.OLGetEdition(idArg)

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

	work, err := OLAPI.OLGetWork(helper.ExtractResourceId(edition.Works[0].ID))

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

func HandleGetBookExistenceSlice(context *gin.Context) {
	bookExistenceSlice, errSlice := helper.FetchBookExistenceSlice()

	if len(errSlice) != 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": errorMessage,
		})

		return
	}

	if len(bookExistenceSlice) == 0 {
		context.Status(http.StatusNoContent)

		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   bookExistenceSlice,
	})
}

func HandleGetBookSearch(context *gin.Context) {
	query := context.Query("query")

	if query == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": fmt.Sprintf("Invalid search term '%s' provided in query parameter 'query'.", query),
		})

		return
	}

	results, err := OLAPI.OLSearchBook(query)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Unable to fetch book metadata and map to supported data structure.",
		})

		return
	}

	if len(results.Results) == 0 {
		context.Status(http.StatusNoContent)

		return
	}

	mappedResults := helper.MapSearchResultSlice(results.Results)

	if len(mappedResults) == 0 {
		context.Status(http.StatusNoContent)

		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   mappedResults,
	})
}

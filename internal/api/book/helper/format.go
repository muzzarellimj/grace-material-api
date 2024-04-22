package helper

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	model "github.com/muzzarellimj/grace-material-api/internal/model/book"
	OLModel "github.com/muzzarellimj/grace-material-api/internal/model/third_party/openlibrary.org"
	"github.com/muzzarellimj/grace-material-api/internal/util"
)

func MapSearchResultSlice(input []OLModel.OLBookSearchResult) []model.BookSearchResult {
	var resultSlice []model.BookSearchResult

	for _, result := range input {
		if len(result.ID) == 0 {
			fmt.Fprint(os.Stdout, "Unable to map OL search result; result did not contain an edition identifier.\n", result)

			continue
		}

		id := result.ID[0]

		if len(result.PublishDate) == 0 {
			fmt.Fprintf(os.Stdout, "Unable to map OL search result; result did not contain a publication date: '%s'\n", id)

			continue
		}

		if len(result.Authors) == 0 {
			fmt.Fprintf(os.Stdout, "Unable to map OL search result; result did not contain an author name set: '%s'\n", id)

			continue
		}

		var publishDate int64

		for _, date := range result.PublishDate {
			datetime := util.ParseDateTime(date)

			if datetime != 0 {
				publishDate = datetime
			}
		}

		if publishDate == 0 {
			fmt.Fprintf(os.Stdout, "Unable to map OL search result; result did not contain a parseable publication date: '%s'\n", id)

			continue
		}

		mappedResult := model.BookSearchResult{
			ID:          id,
			Title:       result.Title,
			Authors:     result.Authors,
			PublishDate: publishDate,
			Image:       FormatImagePath(id),
		}

		resultSlice = append(resultSlice, mappedResult)
	}

	return resultSlice
}

// Format an ISBN-10 or ISBN-13 to remove dashes; e.g., "978-0000000000" becomes "9780000000000".
//
// Return: formatted ISBN when an input string is provided, an empty string when one is not.
func FormatISBN(isbn string) string {
	return strings.ReplaceAll(isbn, "-", "")
}

func FormatImagePath(id string) string {
	return fmt.Sprintf("https://covers.openlibrary.org/b/olid/%s-L.jpg", id)
}

// Extract an edition description, which can either be in string or map[string]string form.
//
// Return: formatted description when one can be parsed, an empty string when one cannot.
func ExtractDescription(value any) string {
	switch t := value.(type) {
	case string:
		return value.(string)
	case map[string]interface{}:
		return value.(map[string]interface{})["value"].(string)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported book description encountered: %v\n", t)
		return ""
	}
}

// Extract an ISBN from a provided slice and format it to remove dashes.
//
// Return: formatted ISBN when an input slice is provided, and empty string when one is not.
func ExtractISBN(slice []string) string {
	if len(slice) == 0 {
		return ""
	}

	return FormatISBN(slice[0])
}

// Extract the first, middle, and last name components from a full name.
//
// Return: first, middle, and last name, each is empty as necessary.
func ExtractName(name string) (string, string, string) {
	nameSlice := strings.Split(name, " ")

	switch length := len(nameSlice); length {

	case 1:
		return nameSlice[0], "", ""

	case 2:
		return nameSlice[0], "", nameSlice[1]

	case 3:
		return nameSlice[0], nameSlice[1], nameSlice[2]

	default:
		return name, "", ""

	}
}

// Extract an OL resource identifier from a resource reference key; e.g., "/books/OL...M" becomes "OL...M".
//
// Return: extracted resource identifier when an input string is provided, an empty string when one is not.
func ExtractResourceId(key string) string {
	pattern := regexp.MustCompile("OL[A-Z0-9]+[A-Z]")

	return pattern.FindString(key)
}

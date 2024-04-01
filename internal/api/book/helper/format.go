package helper

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Format an ISBN-10 or ISBN-13 to remove dashes; e.g., "978-0000000000" becomes "9780000000000".
//
// Return: formatted ISBN when an input string is provided, an empty string when one is not.
func FormatISBN(isbn string) string {
	return strings.ReplaceAll(isbn, "-", "")
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

package helper

import (
	"regexp"
	"strings"
)

// Format an ISBN-10 or ISBN-13 to remove dashes; e.g., "978-0000000000" becomes "9780000000000".
//
// Return: formatted ISBN when an input string is provided, an empty string when one is not.
func FormatISBN(isbn string) string {
	return strings.ReplaceAll(isbn, "-", "")
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

func ExtractResourceId(key string) string {
	pattern := regexp.MustCompile("OL[A-Z0-9]+[A-Z]")

	return pattern.FindString(key)
}

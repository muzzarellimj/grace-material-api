package util

import "strings"

// Format a string value to be execution-ready with PSQL queries by replacing single-quotes and double-quotes with
// single-quote and double-quote pairs, respectively.
//
// Return: formatted PSQL-safe string when an input string is provided, an empty string when one is not.
func FormatPSQLString(value string) string {
	value = strings.ReplaceAll(value, "'", "''")
	value = strings.ReplaceAll(value, "\"", "\"\"")

	return value
}

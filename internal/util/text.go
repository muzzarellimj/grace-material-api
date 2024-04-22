package util

import (
	"strings"
	"time"
)

// Format a date-time string value with a supported layout (format) set to be replaced with a Unix timestamp.
//
// Return: parsed timestamp with sucess, 0 without.
func ParseDateTime(value string) int64 {
	var timestamp int64

	layoutSlice := []string{"January 2, 2006", "Jan 2, 2006", "2006-01-02", "2006/01/02", "01-02-2006", "01/02/2006"}

	for _, layout := range layoutSlice {
		time, err := time.Parse(layout, value)

		if err != nil {
			continue
		}

		timestamp = time.Unix()
	}

	return timestamp
}

// Format a string value to be execution-ready with PSQL queries by replacing single-quotes and double-quotes with
// single-quote and double-quote pairs, respectively.
//
// Return: formatted PSQL-safe string when an input string is provided, an empty string when one is not.
func FormatPSQLString(value string) string {
	value = strings.ReplaceAll(value, "'", "''")
	value = strings.ReplaceAll(value, "\"", "\"\"")

	return value
}

package util_test

import (
	"testing"

	"github.com/muzzarellimj/grace-material-api/pkg/util"
)

func TestParseDateTimeReturnsTimestamp(t *testing.T) {
	value := "November 10, 2022"

	actual := util.ParseDateTime(value)
	expected := int64(1668038400)

	if actual != expected {
		t.Fatalf("Actual timestamp '%d' does not match expected timestamp '%d'.", actual, expected)
	}
}

func TestFormatPSQLStringReturnsSingleQuotes(t *testing.T) {
	value := "Children's fiction"

	actual := util.FormatPSQLString(value)
	expected := "Children''s fiction"

	if actual != expected {
		t.Fatalf("Actual formatted string '%s' does not match expected formatted string '%s'.", actual, expected)
	}
}

func TestFormatPSQLStringReturnsDoubleQuotes(t *testing.T) {
	value := "Russell \"Russ\" Vitale"

	actual := util.FormatPSQLString(value)
	expected := "Russell \"\"Russ\"\" Vitale"

	if actual != expected {
		t.Fatalf("Actual formatted string '%s' does not match expected formatted string '%s'.", actual, expected)
	}
}

func TestFormatPSQLStringReturnsEmptyString(t *testing.T) {
	actual := util.FormatPSQLString("")

	if actual != "" {
		t.Fatalf("Actual formatted string '%s' does not match expected empty formatted string.", actual)
	}
}

package helper_test

import (
	"testing"

	"github.com/muzzarellimj/grace-material-api/pkg/api/book/helper"
)

func TestFormatISBNReturnsISBN10(t *testing.T) {
	isbn := "0140328726"

	expected := "0140328726"
	actual := helper.FormatISBN(isbn)

	if actual != expected {
		t.Fatalf("Actual ISBN-10 '%s' does not match expected ISBN-10 '%s'.", actual, expected)
	}
}

func TestFormatISBNReturnsISBN13(t *testing.T) {
	isbn := "978-0140328721"

	expected := "9780140328721"
	actual := helper.FormatISBN(isbn)

	if actual != expected {
		t.Fatalf("Actual ISBN-13 '%s' does not match expected ISBN-13 '%s'.", actual, expected)
	}
}

func TestFormatISBNReturnsEmptyString(t *testing.T) {
	actual := helper.FormatISBN("")

	if actual != "" {
		t.Fatalf("Actual ISBN '%s' does not match expected empty ISBN.", actual)
	}
}

func TestExtractISBNReturnsISBN10(t *testing.T) {
	slice := []string{"0140328726"}

	expected := "0140328726"
	actual := helper.ExtractISBN(slice)

	if actual != expected {
		t.Fatalf("Actual ISBN-10 '%s' does not match expected ISBN-10 '%s'.", actual, expected)
	}
}

func TestExtractISBNReturnsISBN13(t *testing.T) {
	slice := []string{"978-0140328721"}

	expected := "9780140328721"
	actual := helper.ExtractISBN(slice)

	if actual != expected {
		t.Fatalf("Actual ISBN-13 '%s' does not match expected ISBN-13 '%s'.", actual, expected)
	}
}

func TestExtractISBNReturnsEmptyString(t *testing.T) {
	actual := helper.ExtractISBN([]string{})

	if actual != "" {
		t.Fatalf("Actual ISBN '%s' does not match expected empty ISBN.", actual)
	}
}

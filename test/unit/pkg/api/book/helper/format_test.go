package helper_test

import (
	"testing"

	"github.com/muzzarellimj/grace-material-api/internal/api/book/helper"
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

func TestExtractNameReturnsFirst(t *testing.T) {
	name := "Confucius"

	actualFirst, actualMiddle, actualLast := helper.ExtractName(name)
	expectedFirst := name

	if actualFirst != expectedFirst {
		t.Fatalf("Actual first name '%s' does not match expected first name '%s'.", actualFirst, expectedFirst)
	}

	if actualMiddle != "" || actualLast != "" {
		t.Fatalf("Actual middle name '%s' and last name '%s' does not match expected empty middle and last names.", actualMiddle, actualLast)
	}
}

func TestExtractNameReturnsFirstLast(t *testing.T) {
	name := "Andrzej Sapkowski"

	actualFirst, actualMiddle, actualLast := helper.ExtractName(name)
	expectedFirst := "Andrzej"
	expectedLast := "Sapkowski"

	if actualFirst != expectedFirst {
		t.Fatalf("Actual first name '%s' does not match expected first name '%s'.", actualFirst, expectedFirst)
	}

	if actualMiddle != "" {
		t.Fatalf("Actual middle name '%s' does not match expected empty middle name.", actualMiddle)
	}

	if actualLast != expectedLast {
		t.Fatalf("Actual last name '%s' does not match expected last name '%s'.", actualLast, expectedLast)
	}
}

func TestExtractNameReturnsFirstMiddleLast(t *testing.T) {
	name := "Edgar Allen Poe"

	actualFirst, actualMiddle, actualLast := helper.ExtractName(name)
	expectedFirst := "Edgar"
	expectedMiddle := "Allen"
	expectedLast := "Poe"

	if actualFirst != expectedFirst {
		t.Fatalf("Actual first name '%s' does not match expected first name '%s'.", actualFirst, expectedFirst)
	}

	if actualMiddle != expectedMiddle {
		t.Fatalf("Actual middle name '%s' does not match expected middle name '%s'.", actualMiddle, expectedMiddle)
	}

	if actualLast != expectedLast {
		t.Fatalf("Actual last name '%s' does not match expected last name '%s'.", actualLast, expectedLast)
	}
}

func TestExtractNameReturnsEmpty(t *testing.T) {
	actualFirst, actualMiddle, actualLast := helper.ExtractName("")

	if actualFirst != "" || actualMiddle != "" || actualLast != "" {
		t.Fatalf("Actual first name '%s', middle name '%s', last name '%s' do not match expected empty first, middle, and last name.", actualFirst, actualMiddle, actualLast)
	}
}

func TestExtractResourceIdReturnsResourceId(t *testing.T) {
	key := "/books/OL0M"

	expected := "OL0M"
	actual := helper.ExtractResourceId(key)

	if actual != expected {
		t.Fatalf("Actual resource identifier '%s' does not match expected resource identifier '%s'.", actual, expected)
	}
}

func TestExtractResourceIdReturnsEmptyString(t *testing.T) {
	actual := helper.ExtractResourceId("")

	if actual != "" {
		t.Fatalf("Actual resource identifier '%s' does not match expected empty resource identifier.", actual)
	}
}

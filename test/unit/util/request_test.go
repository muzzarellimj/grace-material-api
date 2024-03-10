package util_test

import (
	"testing"

	"github.com/muzzarellimj/grace-material-api/util"
)

func TestCreateRequestPath(t *testing.T) {
	expected := "https://api.grace.com/endpoint/"

	path, err := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "", make(map[string]string))

	if path != expected || err != nil {
		t.Fatalf("Expected path %s does not match actual path %s: %v\n", expected, path, err)
	}
}

func TestCreateRequestPathRouteParameter(t *testing.T) {
	expected := "https://api.grace.com/endpoint/1"

	path, err := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "1", make(map[string]string))

	if path != expected || err != nil {
		t.Fatalf("Expected path %s does not match actual path %s: %v\n", expected, path, err)
	}
}

func TestCreateRequestPathQueryParameter(t *testing.T) {
	expected := "https://api.grace.com/endpoint/?id=1&"

	queryParams := make(map[string]string)
	queryParams["id"] = "1"

	path, err := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "", queryParams)

	if path != expected || err != nil {
		t.Fatalf("Expected path %s does not match actual path %s: %v\n", expected, path, err)
	}
}

func TestCreateRequest(t *testing.T) {
	expectedMethod := "GET"
	expectedPath := "https://api.grace.com/endpoint/1"

	path, _ := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "1", make(map[string]string))
	request, err := util.CreateRequest("GET", path, make(map[string]string))

	if err != nil {
		t.Fatalf("Unable to create test request: %v\n", err)
	}

	if request.Method != expectedMethod {
		t.Fatalf("Expected request method %s does not match actual request method %s: %v\n", expectedMethod, request.Method, err)
	}

	if request.URL.String() != expectedPath {
		t.Fatalf("Expected request URL %s does not match actual request URL %s: %v\n", expectedPath, request.URL.String(), err)
	}
}

func TestCreateRequestAuthorization(t *testing.T) {
	expectedMethod := "GET"
	expectedPath := "https://api.grace.com/endpoint/1"
	expectedHeader := "Bearer GraceTestToken"

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer GraceTestToken"

	path, _ := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "1", make(map[string]string))
	request, err := util.CreateRequest("GET", path, headers)

	if err != nil {
		t.Fatalf("Unable to create test request: %v\n", err)
	}

	if request.Method != expectedMethod {
		t.Fatalf("Expected request method %s does not match actual request method %s: %v\n", expectedMethod, request.Method, err)
	}

	if request.URL.String() != expectedPath {
		t.Fatalf("Expected request URL %s does not match actual request URL %s: %v\n", expectedPath, request.URL.String(), err)
	}

	if request.Header.Get("Authorization") != expectedHeader {
		t.Fatalf("Expected request header %s does not match actual request header %s: %v\n", expectedHeader, request.Header.Get("Authorization"), err)
	}
}

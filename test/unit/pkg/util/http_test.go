package util_test

import (
	"net/http"
	"testing"

	"github.com/muzzarellimj/grace-material-api/pkg/util"
)

func TestCreateRequestPathReturnsSimplePath(t *testing.T) {
	expected := "https://api.grace.com/endpoint/"

	path, err := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "", make(map[string]string))

	if err != nil {
		t.Fatalf("Unable to create request path: %v\n", err)
	}

	if path != expected {
		t.Fatalf("Actual path '%s' does not match expected path '%s': %v\n", path, expected, err)
	}
}

func TestCreateRequestPathReturnsRouteArgPath(t *testing.T) {
	expected := "https://api.grace.com/endpoint/1"

	path, err := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "1", make(map[string]string))

	if err != nil {
		t.Fatalf("Unable to create request path: %v\n", err)
	}

	if path != expected {
		t.Fatalf("Actual path '%s' does not match expected path '%s': %v\n", path, expected, err)
	}
}

func TestCreateRequestPathReturnsQueryArgPath(t *testing.T) {
	expected := "https://api.grace.com/endpoint/?id=1&"

	path, err := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "", map[string]string{"id": "1"})

	if err != nil {
		t.Fatalf("Unable to create request path: %v\n", err)
	}

	if path != expected {
		t.Fatalf("Actual path '%s' does not match expected path '%s': %v\n", path, expected, err)
	}
}

func TestCreateRequestPathHandlesEmptyBaseArg(t *testing.T) {
	path, err := util.CreateRequestPath("", "", "", make(map[string]string))

	if err == nil {
		t.Fatalf("Unable to catch error with invalid 'base' and 'route' arguments.")
	}

	if path != "" {
		t.Fatalf("Actual path '%s' does not match expected empty path: %v\n", path, err)
	}
}

func TestCreateRequestReturnsSimpleRequest(t *testing.T) {
	expectedMethod := "GET"
	expectedPath := "https://api.grace.com/endpoint/1"

	path, _ := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "1", make(map[string]string))
	request, err := util.CreateRequest(http.MethodGet, path, make(map[string]string))

	if err != nil {
		t.Fatalf("Unable to create request: %v\n", err)
	}

	if request.Method != expectedMethod {
		t.Fatalf("Actual request method '%s' does not match expected request method '%s': %v\n", request.Method, expectedMethod, err)
	}

	if request.URL.String() != expectedPath {
		t.Fatalf("Actual request path '%s' does not match expected request path '%s': %v\n", request.URL.String(), expectedPath, err)
	}
}

func TestCreateRequestReturnsAuthorizedRequest(t *testing.T) {
	expectedMethod := "POST"
	expectedPath := "https://api.grace.com/endpoint/?id=1&"
	expectedAuthorization := "Bearer GraceTestToken"

	path, _ := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "", map[string]string{"id": "1"})
	request, err := util.CreateRequest(http.MethodPost, path, map[string]string{"Authorization": "Bearer GraceTestToken"})

	if err != nil {
		t.Fatalf("Unable to create request: %v\n", err)
	}

	if request.Method != expectedMethod {
		t.Fatalf("Actual request method '%s' does not match expected request method '%s': %v\n", request.Method, expectedMethod, err)
	}

	if request.URL.String() != expectedPath {
		t.Fatalf("Actual request path '%s' does not match expected request path '%s': %v\n", request.URL.String(), expectedPath, err)
	}

	if request.Header.Get("Authorization") != expectedAuthorization {
		t.Fatalf("Actual authorization header '%s' does not match expected authorization header '%s': %v\n", request.URL.String(), expectedPath, err)
	}
}

func TestCreateRequestHandlesEmptyMethodArg(t *testing.T) {
	path, _ := util.CreateRequestPath("https://api.grace.com", "/endpoint/", "", make(map[string]string))
	request, err := util.CreateRequest("", path, make(map[string]string))

	if err == nil {
		t.Fatal("Unable to catch error with empty 'method' argument.\n")
	}

	if request != nil {
		t.Fatalf("Actual request '%v' does not match expected nil request: %v\n", request, err)
	}
}

func TestCreateRequestHandlesEmptyPathArg(t *testing.T) {
	request, err := util.CreateRequest(http.MethodGet, "", make(map[string]string))

	if err == nil {
		t.Fatal("Unable to catch error with empty 'path' argument.\n")
	}

	if request != nil {
		t.Fatalf("Actual request '%v' does not match expected nil request: %v\n", request, err)
	}
}

func TestExecuteRequestReturnsStatusOk(t *testing.T) {
	path, _ := util.CreateRequestPath("https://google.com", "", "", make(map[string]string))
	request, _ := util.CreateRequest(http.MethodGet, path, make(map[string]string))
	response, err := util.ExecuteRequest(request)

	if err != nil {
		t.Fatalf("Unable to execute request: %v\n", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("Actual request '%v' does not match expected nil request: %v\n", request, err)
	}
}
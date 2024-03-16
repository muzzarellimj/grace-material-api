package util

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Create an HTTP request path with a base URL and route, and optional route parameter and query parameter map.
//
// Return: built request path and nil with success, empty string and error without.
func CreateRequestPath(base string, route string, routeParam string, queryParams map[string]string) (string, error) {
	var builder strings.Builder

	if base == "" {
		err := errors.New("invalid 'base' argument provided")

		fmt.Fprintf(os.Stderr, "Unable to create request path without required 'base' argument: %v\n", err)

		return "", err
	}

	_, err := builder.WriteString(fmt.Sprint(base, route))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request path: %v\n", err)

		return "", err
	}

	if routeParam != "" {
		_, err := builder.WriteString(fmt.Sprint("/", routeParam))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create request path with route parameter: %v\n", err)

			return "", err
		}
	}

	if len(queryParams) > 0 {
		builder.WriteString("?")

		for key, value := range queryParams {
			_, err := builder.WriteString(fmt.Sprint(key, "=", value, "&"))

			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to create request path with query parameter(s): %v\n", err)

				return "", err
			}
		}
	}

	path := builder.String()
	path = strings.ReplaceAll(path, " ", "%20")
	path = strings.TrimSuffix(path, "&")

	return path, nil
}

// Create an HTTP request with a request method, request path, and header map.
//
// Return: built request and nil with success, nil and error without.
func CreateRequest(method string, path string, headers map[string]string) (*http.Request, error) {
	if method == "" || path == "" {
		err := errors.New("invalid 'method' or 'path' argument provided")

		fmt.Fprintf(os.Stderr, "Unable to create request without required 'method' and 'path' arguments: %v\n", err)

		return nil, err
	}

	request, err := http.NewRequest(method, path, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create HTTP request: %v\n", err)

		return nil, err
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	return request, nil
}

// Execute an HTTP request on a new client.
//
// Return: response and nil with success, nil and error without.
func ExecuteRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to execute HTTP request: %v\n", err)

		return &http.Response{}, err
	}

	return response, nil
}

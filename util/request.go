package util

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Create an HTTP request path with a base URL, endpoint route, route parameter, and query parameter map.
//
// Return: built request path and nil with success, empty string and error without.
func CreateRequestPath(base string, endpoint string, routeParam string, queryParams map[string]string) (string, error) {
	var builder strings.Builder

	_, err := builder.WriteString(fmt.Sprint(base, endpoint))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request path: %v\n", err)

		return "", err
	}

	if routeParam != "" {
		_, err := builder.WriteString(routeParam)

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

	return builder.String(), nil
}

// Create an HTTP request with an HTTP method, a request path, and header map.
//
// Return: request and nil with success, nil and error without.
func CreateRequest(method string, path string, headers map[string]string) (*http.Request, error) {
	request, err := http.NewRequest(method, path, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create request path with query parameter(s): %v\n", err)

		return nil, err
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	return request, nil
}

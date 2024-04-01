package api_test

import (
	"testing"

	api "github.com/muzzarellimj/grace-material-api/internal/api/third_party/twitch.tv"
)

// Note: needs to be run with injected client_id and client_secret - otherwise will return 400.
func TestTTVAuthenticateRequestReturnsToken(t *testing.T) {
	token, err := api.TTVAuthenticateRequest()

	if err != nil {
		t.Fatalf("Unable to authenticate request and parse token in response: %v\n", err)
	}

	if token == "" {
		t.Fatalf("Actual empty token does not match expected non-empty token.")
	}
}

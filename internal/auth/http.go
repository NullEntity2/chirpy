package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authz := headers.Get("Authorization")
	if authz == "" {
		return "", fmt.Errorf("headers does not include authorization header")
	}
	trimmedAuthz, _ := strings.CutPrefix(authz, "Bearer ")
	return trimmedAuthz, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authz := headers.Get("Authorization")
	if authz == "" {
		return "", fmt.Errorf("headers does not include authorization header")
	}
	trimmedAuthz, _ := strings.CutPrefix(authz, "ApiKey ")
	return trimmedAuthz, nil
}

package auth

import (
	"net/http"
	"testing"
)

func TestBearerHappy(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer helloworld")
	token, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("error getting bearer token: %s", err)
	}

	if token != "helloworld" {
		t.Fatalf("token expected to be helloworld, got %s", token)
	}
}

func TestBearerNoAuthzHeader(t *testing.T) {
	header := http.Header{}
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatalf("did not recieve expected error when calling GetBearerToken")
	}
}

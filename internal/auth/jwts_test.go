package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	const secret = "secret"
	uuid := uuid.MustParse("2a8e1bba-457a-47d7-9501-e0dd9cc44bba")
	jwt, err := MakeJWT(uuid, secret, 5*time.Minute)
	if err != nil {
		t.Errorf("error creating jwt: %s", err)
	}

	parsedUuid, err := ValidateJWT(jwt, secret)
	if err != nil {
		t.Errorf("error validating jwt: %s", err)
	}

	if uuid != parsedUuid {
		t.Errorf("uuid != parsedUuid")
	}
}

func TestValidateExpiredJWT(t *testing.T) {
	const secret = "secret"
	uuid := uuid.MustParse("2a8e1bba-457a-47d7-9501-e0dd9cc44bba")
	token, err := MakeJWT(uuid, secret, -5*time.Minute)
	if err != nil {
		t.Errorf("error creating jwt: %s", err)
	}

	_, err = ValidateJWT(token, secret)
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("error validating jwt: %s", err)
	}
}

func TestValidateWrongSecret(t *testing.T) {
	const secret = "secret"
	uuid := uuid.MustParse("2a8e1bba-457a-47d7-9501-e0dd9cc44bba")
	token, err := MakeJWT(uuid, "different secret", 5*time.Minute)
	if err != nil {
		t.Errorf("error creating jwt: %s", err)
	}

	_, err = ValidateJWT(token, secret)
	if err != nil && !errors.Is(err, jwt.ErrSignatureInvalid) {
		t.Errorf("error validating jwt: %s", err)
	}
}

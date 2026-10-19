package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "Hello World"
	hashed_password, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password: %s", err)
	}
	if password == hashed_password {
		t.Errorf("Error password did not hash correctly")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "Hello World"
	hashed_password, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password: %s", err)
	}

	valid, err := CheckPasswordHash(password, hashed_password)
	if err != nil {
		t.Errorf("error checking password hash: %s", err)
	}

	if !valid {
		t.Errorf("Hashed password is not valid")
	}
}

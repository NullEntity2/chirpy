package main

import "testing"

func TestSanitizeString(t *testing.T) {
	result := sanitizeString("This is a kerfuffle opinion I need to share with the world")
	if result != "This is a **** opinion I need to share with the world" {
		t.Errorf(`Expected string to be sanitized got %s`, result)
	}
}

func TestSanitizeString2(t *testing.T) {
	result := sanitizeString("I hear Mastodon is better than Chirpy. sharbert I need to migrate")
	if result != "I hear Mastodon is better than Chirpy. **** I need to migrate" {
		t.Errorf(`Expected string to be sanitized got %s`, result)
	}
}

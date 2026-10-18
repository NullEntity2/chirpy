package main

import (
	"strings"
)

var bannedWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func sanitizeString(input string) string {
	inputs := strings.Split(input, " ")
	for k, v := range inputs {
		lowerV := strings.ToLower(v)
		if _, ok := bannedWords[lowerV]; ok {
			inputs[k] = "****"
		}
	}
	return strings.Join(inputs, " ")
}

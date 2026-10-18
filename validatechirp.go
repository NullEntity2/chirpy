package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateChirp(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type successStruct struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		resp.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(resp, 400, "Chirp is too long")
	} else {
		suc := successStruct{CleanedBody: sanitizeString(params.Body)}
		respondWithJSON(resp, 200, suc)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorStruct struct {
		Error string `json:"error"`
	}
	error := errorStruct{Error: msg}
	respondWithJSON(w, code, error)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

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

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NullEntity2/chirpy/internal/auth"
	"github.com/NullEntity2/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) postPolkaWebookFunc(resp http.ResponseWriter, req *http.Request) {
	type polkaRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Printf("Error getting API key from headers: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	if apiKey != cfg.PolkaKey {
		log.Printf("apiKey != cfg.PolkaKey")
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := polkaRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	if params.Event != "user.upgraded" {
		log.Printf("Got webhook event: %s", params.Event)
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Error parsing user id: %s", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = cfg.dbQueries.GetUser(req.Context(), userId)
	if err != nil {
		log.Printf("User not found: %s", err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	chirpyRedParams := database.SetUserChirpyRedParams{
		ID:          userId,
		IsChirpyRed: true,
	}
	err = cfg.dbQueries.SetUserChirpyRed(req.Context(), chirpyRedParams)
	if err != nil {
		log.Printf("Error setting user chirpy red: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}

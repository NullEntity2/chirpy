package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/NullEntity2/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) postChirpsFunc(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := request{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Message must be less than 141 characters")
		return
	}

	ccParams := database.CreateChirpParams{
		Body:   sanitizeString(params.Body),
		UserID: params.UserID,
	}
	createdChirp, err := cfg.dbQueries.CreateChirp(r.Context(), ccParams)
	if err != nil {
		log.Printf("Error inserting user: %s", err)
		w.WriteHeader(500)
		return
	}

	resp := response{
		Id:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, resp)
}

func (cfg *apiConfig) getChirpsFunc(resp http.ResponseWriter, req *http.Request) {
	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirps, err := cfg.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		resp.WriteHeader(500)
		return
	}

	responses := make([]response, 0, len(chirps))
	for _, chirp := range chirps {
		r := response{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		responses = append(responses, r)
	}
	respondWithJSON(resp, http.StatusOK, responses)
}

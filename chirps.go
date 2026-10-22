package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/NullEntity2/chirpy/internal/auth"
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

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Did not recieve bearer token: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uuid, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("Could not validate bearer token: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Message must be less than 141 characters")
		return
	}

	ccParams := database.CreateChirpParams{
		Body:   sanitizeString(params.Body),
		UserID: uuid,
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

type chirpResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirpsFunc(resp http.ResponseWriter, req *http.Request) {

	chirps, err := cfg.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		resp.WriteHeader(500)
		return
	}

	responses := make([]chirpResponse, 0, len(chirps))
	for _, chirp := range chirps {
		r := chirpResponse{
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

func (cfg *apiConfig) getChirpFunc(resp http.ResponseWriter, req *http.Request) {
	idString := req.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Printf("Error parsing uuid from idString: %s", err)
		respondWithError(resp, http.StatusBadRequest, "error parsing provided uuid")
		return
	}

	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Chirp not found: %s", err)
			resp.WriteHeader(404)
			return
		} else {
			log.Printf("Error getting chirps: %s", err)
			resp.WriteHeader(500)
			return
		}
	}

	r := chirpResponse{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(resp, http.StatusOK, r)
}

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/NullEntity2/chirpy/internal/auth"
	"github.com/NullEntity2/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

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

func (cfg *apiConfig) getChirpsFunc(resp http.ResponseWriter, req *http.Request) {
	authorIdString := req.URL.Query().Get("author_id")
	authorId, _ := uuid.Parse(authorIdString)

	sort := req.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error
	if authorId != uuid.Nil {
		chirps, err = cfg.dbQueries.GetAllChirpsByAuthor(req.Context(), authorId)
		if err != nil {
			log.Printf("Error getting chirps by author: %s", err)
			resp.WriteHeader(500)
			return
		}
	} else {
		chirps, err = cfg.dbQueries.GetAllChirps(req.Context())
		if err != nil {
			log.Printf("Error getting chirps: %s", err)
			resp.WriteHeader(500)
			return
		}
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

	switch sort {
	case "asc":
		slices.SortFunc(responses, func(a, b chirpResponse) int {
			return a.CreatedAt.Compare(b.CreatedAt)
		})
	case "desc":
		slices.SortFunc(responses, func(a, b chirpResponse) int {
			return b.CreatedAt.Compare(a.CreatedAt)
		})
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

func (cfg *apiConfig) deleteChirpFunc(resp http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Did not recieve bearer token: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwtUserId, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		log.Printf("Could not validate bearer token: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirpIdString := req.PathValue("chirpID")
	chirpId, err := uuid.Parse(chirpIdString)
	if err != nil {
		log.Printf("Error parsing uuid from chirpIdString: %s", err)
		respondWithError(resp, http.StatusBadRequest, "error parsing provided uuid")
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Could not find chirp: %s", err)
			respondWithError(resp, http.StatusNotFound, "error parsing provided uuid")
			return
		} else {
			log.Printf("Error querying chirp: %s", err)
			respondWithError(resp, http.StatusInternalServerError, "error parsing provided uuid")
			return
		}
	}

	if jwtUserId != chirp.UserID {
		log.Printf("Chirp does not belong to owner")
		respondWithError(resp, http.StatusForbidden, "error parsing provided uuid")
		return
	}

	err = cfg.dbQueries.DeleteChirp(req.Context(), chirp.ID)
	if err != nil {
		log.Printf("Error deleting chirp: %s", err)
		respondWithError(resp, http.StatusInternalServerError, "error parsing provided uuid")
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}

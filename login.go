package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/NullEntity2/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) postLoginFunc(resp http.ResponseWriter, req *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := request{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		resp.WriteHeader(500)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error fetching user from db: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		resp.Write([]byte("Incorrect email or password"))
		return
	}

	valid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password hash: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		resp.Write([]byte("Incorrect email or password"))
		return
	}

	if !valid {
		resp.WriteHeader(http.StatusUnauthorized)
		resp.Write([]byte("Incorrect email or password"))
		return
	}

	r := response{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(resp, http.StatusOK, r)
}

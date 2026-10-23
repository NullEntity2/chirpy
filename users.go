package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/NullEntity2/chirpy/internal/auth"
	"github.com/NullEntity2/chirpy/internal/database"
	"github.com/google/uuid"
)

type userResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) postUsersFunc(resp http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := userRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		resp.WriteHeader(500)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
	}

	createUserParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	createdUser, err := cfg.dbQueries.CreateUser(req.Context(), createUserParams)
	if err != nil {
		log.Printf("Error inserting user: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	respStruct := userResponse{
		Id:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}
	respondWithJSON(resp, http.StatusCreated, respStruct)
}

func (cfg *apiConfig) putUsersFunc(resp http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error checking password hash: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.JWTSecret)
	if err != nil {
		log.Printf("Failed to validate JWT: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := userRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		resp.WriteHeader(500)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("error hashing password: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	uup := database.UpdateUserParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	err = cfg.dbQueries.UpdateUser(req.Context(), uup)
	if err != nil {
		log.Printf("error updating user: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), uup.Email)

	ur := userResponse{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(resp, http.StatusOK, ur)
}

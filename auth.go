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

func (cfg *apiConfig) postLoginFunc(resp http.ResponseWriter, req *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Id           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
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

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Duration(1)*time.Hour)
	if err != nil {
		log.Printf("Error making JWT: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	crtParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(60*24) * time.Hour),
	}
	dbToken, err := cfg.dbQueries.CreateRefreshToken(req.Context(), crtParams)

	r := response{
		Id:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: dbToken.Token,
		IsChirpyRed:  user.IsChirpyRed,
	}
	respondWithJSON(resp, http.StatusOK, r)
}

func (cfg *apiConfig) postRefreshFunc(resp http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error checking password hash: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = auth.ValidateJWT(bearerToken, cfg.JWTSecret)
	if err != nil {
		log.Printf("Failed to validate JWT: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), bearerToken)
	if err != nil {
		log.Printf("error getting refresh token from db: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	if dbToken.RevokedAt.Valid {
		log.Printf("error token is revoked")
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !dbToken.ExpiresAt.After(time.Now()) {
		log.Printf("error token is expired")
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.MakeJWT(dbToken.UserID, cfg.JWTSecret, time.Duration(1)*time.Hour)
	if err != nil {
		log.Printf("Error making JWT: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	r := response{
		Token: accessToken,
	}
	respondWithJSON(resp, http.StatusOK, r)
}

func (cfg *apiConfig) postRevokeFunc(resp http.ResponseWriter, req *http.Request) {
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error checking password hash: %s", err)
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(req.Context(), bearerToken)
	if err != nil {
		log.Printf("error revoking token: %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}

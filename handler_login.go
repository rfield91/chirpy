package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/rfield91/chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error retrieving user: %s", err)
			respondWithError(w, http.StatusUnauthorized, "Error retrieving user", err)
			return
		} else {
			log.Printf("Error retrieving user: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Error retrieving user", err)
			return
		}
	}

	isPasswordMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)

	if err != nil {
		log.Printf("Error checking password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error checking password", err)
		return
	}

	if !isPasswordMatch {
		log.Printf("Invalid password: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid password", err)
		return
	}

	var expiresIn time.Duration

	if params.ExpiresInSeconds <= 0  || params.ExpiresInSeconds > 3600 {
		expiresIn = time.Duration(1 * time.Hour)
	} else {
		expiresIn = time.Duration(params.ExpiresInSeconds)
	}

	token, err := auth.MakeJWT(user.ID, cfg.authSecret, expiresIn)

	log.Printf("Token: %s", token)
	log.Printf("Expire: %v", expiresIn)
	
	if err != nil {
		log.Printf("Error making JWT: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error making JWT", err)
		return
	}
	
	userJson := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, userJson)
}
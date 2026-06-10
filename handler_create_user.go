package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Invalid request body", err)
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), params.Email)

	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
	}

	userJson := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusCreated, userJson)
}
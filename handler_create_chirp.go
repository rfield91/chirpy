package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/rfield91/chirpy/internal/auth"
	"github.com/rfield91/chirpy/internal/database"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, tokenErr := auth.GetBearerToken(req.Header)

	if tokenErr != nil {
		log.Printf("Could not get auth token: %s", tokenErr)
		respondWithError(w, http.StatusUnauthorized, "Invalid request body", tokenErr)
		return
	}

	userId, jwtValidationErr := auth.ValidateJWT(token, cfg.authSecret)

	if jwtValidationErr != nil {
		log.Printf("Could not validate auth: %s", jwtValidationErr)
		respondWithError(w, http.StatusUnauthorized, "Invalid request body", jwtValidationErr)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too longer", nil)
		return
	}

	body := cleanChirp(params.Body)

	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: body,
		UserID: userId,
	})

	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to create chirp", err)
	}

	chirpJson := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	log.Printf("Chirp: %v", chirpJson)

	respondWithJSON(w, http.StatusCreated, chirpJson)
}

func isBadWord(word string, checkWords []string) bool {
	for _, bad := range checkWords {
		if strings.ToUpper(word) == strings.ToUpper(bad) {
			return true
		} 
	}

	return false
}

func cleanChirp(chirpBody string) string {
	var words []string
	badWords := [...]string{"kerfuffle", "sharbert", "fornax"}

	for _, word := range strings.Split(chirpBody, " ") {
		log.Printf("word: %s", word)

		if isBadWord(word, badWords[:]) == true {
			words = append(words, "****")
		} else {
			words = append(words, word)
		}
	}

	cleanedBody := strings.Join(words, " ")

	return cleanedBody
}
package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, req *http.Request) {
	chirpId := req.PathValue("chirpID")

	if len(chirpId) == 0 {
		log.Printf("ChirpID cannot be empty")
		respondWithError(w, http.StatusBadRequest, "ChirpID cannot be empty", nil)
		return
	}

	parsedchirpId, err := uuid.Parse(chirpId)

	if err != nil {
		log.Printf("ChirpID could not be parsed")
		respondWithError(w, http.StatusBadRequest, "ChirpID could not be parsed", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirpById(req.Context(), parsedchirpId)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Chirp not found: %s", err)
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		} else {
			log.Printf("Error retrieving chirps: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Unable to retrieve chirps", err)
			return
		}
	}

	chirpJson := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirpJson)
}
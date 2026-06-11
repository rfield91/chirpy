package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(req.Context())

	if err != nil {
		log.Printf("Error retrieving chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to retrieve chirps", err)
		return
	}

	var chirpsJsonList []Chirp

	for _, chirp := range chirps {
		chirpsJsonList = append(chirpsJsonList, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsJsonList)
}
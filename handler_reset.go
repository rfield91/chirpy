package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
	}

	err := cfg.dbQueries.DeleteUsers(req.Context())

	if err != nil {
		log.Printf("Error deleting users: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to delete users", err)
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)



func isBadWord(word string, checkWords []string) bool {
	for _, bad := range checkWords {
		if strings.ToUpper(word) == strings.ToUpper(bad) {
			log.Printf("Bad word: %s", bad)
			return true
		} 
	}

	return false
}

func handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Invalid request body", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too longer", nil)
		return
	}

	var words []string
	badWords := [...]string{"kerfuffle", "sharbert", "fornax"}

	for _, word := range strings.Split(params.Body, " ") {
		log.Printf("word: %s", word)

		if isBadWord(word, badWords[:]) == true {
			words = append(words, "****")
		} else {
			words = append(words, word)
		}
	}

	cleanedBody := strings.Join(words, " ")

	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: cleanedBody})
}
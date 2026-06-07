package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, req)
		cfg.fileserverHits.Add(1)
		log.Printf("Hits: %v", cfg.fileserverHits.Load())
	})
}
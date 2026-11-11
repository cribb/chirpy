package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handleOK(w http.ResponseWriter, req *http.Request) {
	respondWithOK(w)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, r)
		},
	)
}

func (cfg *apiConfig) getHits(w http.ResponseWriter, req *http.Request) {
	text := fmt.Sprintf(
		`<html> 
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
		</html>`,
		cfg.fileserverHits.Load())

	// text := fmt.Sprintf("Hits: %v\n", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))

}

func (cfg *apiConfig) resetAPI(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Verboten: wrong platform")
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.ResetUsers(req.Context())
	if err != nil {
		fmt.Printf("Error on User reset: %v", err)
		return
	}
	text := fmt.Sprintf("Database is reset.\nReset hits to: %v\n\n", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

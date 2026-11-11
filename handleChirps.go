package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cribb/chirpy/internal/auth"
	"github.com/cribb/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlePostChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// fmt.Printf("params: %v\n", params)
		// fmt.Printf("req.Body: %v\n", req.Body)
		// fmt.Printf("decoder: %v\n", decoder)
		respondWithError(w, 500, fmt.Sprintf("Error decoding JSON: %s", err))
		return
	}
	// fmt.Printf("decoder: %v, params: %v", decoder, params.Body)

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Bad token.")
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid UUID")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	chirp_params := database.CreateChirpParams{
		UserID: userID,
		Body:   badWordsBeGone(params.Body),
	}

	// fmt.Printf("chirp_params: %v\n", chirp_params)

	res, err := cfg.dbQueries.CreateChirp(req.Context(), chirp_params)
	if err != nil {
		respondWithError(w, 400, "Error creating chirp.")
		return
	}

	newChirp := Chirp{
		ID:        res.ID,
		UserID:    res.UserID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Body:      res.Body,
	}

	// fmt.Printf("chirp: %+v\n", res)

	respondWithJSON(w, 201, newChirp)

}

func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, req *http.Request) {

	var chirps []Chirp

	res, err := cfg.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		fmt.Printf("Error occurred while retrieving chirps: %v\n", err)
		return
	}

	for i, c := range res {
		chirps = append(chirps,
			Chirp{
				ID:        c.ID,
				UserID:    c.UserID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
			},
		)
		fmt.Printf("Appending entry: %v\n", i)

	}

	respondWithJSON(w, 200, chirps)

}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, req *http.Request) {

	var chirp Chirp

	// res, err := cfg.dbQueries.GetChirp(req.Context(), chirp_id)

	id_str := req.PathValue("chirpID")

	chirp_id, err := uuid.Parse(id_str)
	if err != nil {
		fmt.Printf("Error processing chirp UUID: %v\n", err)
		return
	}

	c, err := cfg.dbQueries.GetChirp(req.Context(), chirp_id)
	if err != nil {
		fmt.Printf("Error occurred while retrieving chirp (%v): %v\n", chirp_id, err)
		w.WriteHeader(404)
		return
	}

	chirp = Chirp{
		ID:        c.ID,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
	}

	respondWithJSON(w, 200, chirp)

}

func (cfg *apiConfig) handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	// fmt.Printf("decoder: %v, params: %v", decoder, params.Body)

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	respondWithCleanVersion(w, 200, params.Body)

}

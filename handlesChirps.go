package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
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

// func (cfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, req *http.Request) {

// 	var chirps []Chirp

// 	res, err := cfg.dbQueries.GetAllChirps(req.Context())
// 	if err != nil {
// 		fmt.Printf("Error occurred while retrieving chirps: %v\n", err)
// 		return
// 	}

// 	for _, c := range res {
// 		chirps = append(chirps,
// 			Chirp{
// 				ID:        c.ID,
// 				UserID:    c.UserID,
// 				CreatedAt: c.CreatedAt,
// 				UpdatedAt: c.UpdatedAt,
// 				Body:      c.Body,
// 			},
// 		)
// 		// fmt.Printf("Appending entry: %v\n", i)

// 	}

// 	respondWithJSON(w, 200, chirps)

// }

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, req *http.Request) {

	sortAscending := true
	res := []database.Chirp{}
	chirps := []Chirp{}
	var err error

	userUuidString := req.URL.Query().Get("author_id")
	fmt.Printf("userUUIDstring: %+v\n", userUuidString)
	desiredOrder := req.URL.Query().Get("sort")
	fmt.Printf("desiredOrder: %+v\n", desiredOrder)

	if desiredOrder == "desc" {
		sortAscending = false
	}

	if len(userUuidString) == 0 {
		res, err = cfg.dbQueries.GetAllChirps(req.Context())
		if err != nil {
			fmt.Printf("Error occurred while retrieving chirps: %v\n", err)
			return
		}
	} else {
		userUUID, err := uuid.Parse(userUuidString)
		if err != nil {
			respondWithError(w, 400, "Bad UUID.")
		}
		fmt.Printf("userUUID: %+v\n", userUUID)

		res, err = cfg.dbQueries.GetChirpsForUser(req.Context(), userUUID)
		if err != nil {
			fmt.Printf("Error occurred while retrieving chirps: %v\n", err)
			return
		}
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

	if !sortAscending {
		slices.Reverse(chirps)
	}

	respondWithJSON(w, 200, chirps)

}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, req *http.Request) {

	var chirp Chirp

	id_str := req.PathValue("chirpID")

	chirp_id, err := uuid.Parse(id_str)
	if err != nil {
		fmt.Printf("Error processing chirp UUID: %v\n", err)
		return
	}

	c, err := cfg.dbQueries.GetChirp(req.Context(), chirp_id)
	if err != nil {
		fmt.Printf("Error occurred while retrieving chirp (%v): %v\n", chirp_id, err)
		respondWithError(w, 404, "Chirp not found.")
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

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, req *http.Request) {

	jwToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Printf("DeleteChirp: Unauthenticated user.\n %v", err)
		respondWithError(w, 401, "Unauthenticated user.")
		return
	}

	user_id, err := auth.ValidateJWT(jwToken, cfg.jwtSecret)
	if err != nil {
		fmt.Printf("User not found. (%v)\n", err)
		respondWithError(w, 401, "User not found.")
		return
	}

	chirp_id, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		fmt.Printf("Error processing chirp UUID: %v\n", err)
		respondWithError(w, 500, "Malformed chirp UUID.")
		return
	}
	fmt.Printf("Delete chirp: %v\nwith user_id: %v\n", chirp_id, user_id)

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirp_id)
	if err != nil {
		fmt.Printf("Error occurred while retrieving chirp (%v): %v\n", chirp_id, err)
		w.WriteHeader(404)
		return
	}
	fmt.Printf("Found chirp: %+v\n", chirp)

	if !(chirp.ID == chirp_id && chirp.UserID == user_id) {
		respondWithError(w, 403, "User is not chirp author.")
		return
	}
	fmt.Printf("User (%v) is chirp (%v) author.\n", chirp.ID, chirp.UserID)

	chirp_params := database.DeleteChirpParams{
		ID:     chirp_id,
		UserID: user_id,
	}

	err = cfg.dbQueries.DeleteChirp(req.Context(), chirp_params)
	if err != nil {
		fmt.Printf("Chirp not found. (%v):\n   %v\n", chirp_id, err)
		respondWithError(w, 404, "Chirp not found.")
		return
	}

	respondWithJSON(w, 204, "Chirp deleted.")
}

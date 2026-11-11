package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cribb/chirpy/internal/auth"
	"github.com/cribb/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handleAddUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding JSON: %s", err))
		return
	}

	email := params.Email
	hashed_pword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 401, "Error generating hash.")
		return
	}

	user_params := database.CreateUserParams{
		Email:          email,
		HashedPassword: hashed_pword,
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), user_params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating user: %s", err))
		return
	}

	fmt.Printf("user: %v\n", user)
	newUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, 201, newUser)

}

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding JSON: %s", err))
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		fmt.Printf("Error: User not found. (%v)\n", err)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	success, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		fmt.Printf("Error with password check: (%v)\n", err)
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	if !success {
		fmt.Printf("Error: Wrong password for user.\n")
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	expiry := time.Hour
	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiry)
	if err != nil {
		respondWithError(w, 500, "Bad Token/Secret.")
		return
	}

	keySizeBytes := 32
	r_tok := auth.MakeRefreshToken(keySizeBytes)
	r_tokNow := time.Now()

	tokenParams := database.CreateRefreshTokenParams{
		Token:     r_tok,
		UserID:    user.ID,
		CreatedAt: r_tokNow,
		UpdatedAt: r_tokNow,
		ExpiresAt: r_tokNow.Add(time.Hour * 24 * 60),
	}

	err = cfg.dbQueries.CreateRefreshToken(req.Context(), tokenParams)
	if err != nil {
		respondWithError(w, 501, "Refresh token not stored properly.")
	}

	loginUser := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        jwtToken,
		RefreshToken: r_tok,
	}

	// fmt.Printf("loginUser: %+v\n", loginUser)

	respondWithJSON(w, 200, loginUser)
}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, req *http.Request) {

	btoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Bad token.")
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(req.Context(), btoken)
	if err != nil {
		respondWithError(w, 500, "Bad db update.")
		return
	}

	respondWithJSON(w, 204, "")

}

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, req *http.Request) {

	btoken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "Bad token.")
		return
	}

	user, err := cfg.dbQueries.GetUserbyRefreshToken(req.Context(), btoken)
	if err != nil {
		respondWithError(w, 401, "Bad user.")
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, 401, "Bad access token.")
		return
	}

	type tokenResponse struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, 200, tokenResponse{Token: newToken})
}

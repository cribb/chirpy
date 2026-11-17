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

	// fmt.Printf("user: %v\n", user)
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

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		fmt.Printf("Error with password check: (%v)\n", err)
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
		fmt.Printf("Error getting bearer token. (%v)\n", err)
		respondWithError(w, 401, "Bad token.")
		return
	}

	fmt.Printf("btoken: %v\n", btoken)
	err = cfg.dbQueries.RevokeRefreshToken(req.Context(), btoken)
	if err != nil {
		fmt.Printf("Bad revocation: %v\n", err)
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

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, req *http.Request) {

	jwToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		// fmt.Printf("UpdateUser: Bad token. (%v)\n", err)
		respondWithError(w, 401, "Bad token.")
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	userNewParams := parameters{}
	err = decoder.Decode(&userNewParams)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding JSON: %s", err))
		return
	}

	uuid, err := auth.ValidateJWT(jwToken, cfg.jwtSecret)
	if err != nil {
		// fmt.Printf("UpdateUser: User not found. (%v)\n", err)
		respondWithError(w, 401, "User not found.")
		return
	}
	// fmt.Printf("UpdateUser: uuid: (%v)\n", uuid)

	// user, err := cfg.dbQueries.GetUserByUUID(req.Context(), uuid)
	// if err != nil {
	// 	fmt.Printf("UpdateUser: User not found. (%v)\n", err)
	// 	respondWithError(w, 500, "User not found.")
	// 	return
	// }

	// fmt.Printf("UpdateUser: updating from old email: %v\n", user.Email)
	// fmt.Printf("UpdateUser: updating  to  new email: %v\n", userNewParams.Email)
	// fmt.Printf("UpdateUser: updating  to  new pword: %v\n", userNewParams.Password)

	hash, err := auth.HashPassword(userNewParams.Password)
	if err != nil {
		respondWithError(w, 500, "Error hashing password.")
		return
	}

	p := database.UpdateUserParams{
		ID:             uuid,
		Email:          userNewParams.Email,
		HashedPassword: hash,
	}

	err = cfg.dbQueries.UpdateUser(req.Context(), p)
	if err != nil {
		respondWithError(w, 500, "Error updating user.")
		return
	}

	// userchk, err := cfg.dbQueries.GetUserByEmail(req.Context(), userNewParams.Email)
	// if err != nil {
	// 	respondWithError(w, 401, "Error retrieving new email for user.")
	// 	return
	// }

	// fmt.Printf("Updated user's email to: %v\n", p.Email)
	// fmt.Printf("Updated user's email to: %v\n", userchk.Email)

	type userResponse struct {
		Email string `json:"email"`
	}

	respondWithJSON(w, 200, userResponse{Email: userNewParams.Email})
}

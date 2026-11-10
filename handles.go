package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
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

	r_tok := auth.MakeRefreshToken(32)
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

func handleOK(w http.ResponseWriter, req *http.Request) {
	text := "OK"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	// fmt.Printf("Payload: %d\n", payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func badWordsBeGone(msg string) string {
	words := strings.Split(msg, " ")

	for i, word := range words {
		lword := strings.ToLower(word)
		if lword == "kerfuffle" || lword == "sharbert" || lword == "fornax" {
			words[i] = "****" // change word to *** in the __original slice__
		}
	}
	jmsg := strings.Join(words, " ")
	return jmsg
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
		w.WriteHeader(403) // Forbidden
		// w.Write([]byte("Forbidden\n"))
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

func validateChirp(w http.ResponseWriter, req *http.Request) {
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

func respondWithCleanVersion(w http.ResponseWriter, code int, msg string) {
	type cleanResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	cmsg := badWordsBeGone(msg)
	respondWithJSON(w, code, cleanResponse{CleanedBody: cmsg})
}

// func respondWithSuccess(w http.ResponseWriter, code int) {
// 	type validResponse struct {
// 		Valid bool `json:"valid"`
// 	}
// 	respondWithJSON(w, code, validResponse{Valid: true})
// }

// func (cfg *apiConfig) middlewareMetricsReport(next http.Handler) http.Handler {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			fmt.Printf("Hits: %v\n", cfg.fileserverHits.Load())
// 			next.ServeHTTP(w, r)
// 		},
// 	)
// }

// func middlewareLog(next http.Handler) http.Handler {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			log.Printf("%s %s", r.Method, r.URL.Path)
// 			next.ServeHTTP(w, r)
// 		},
// 	)
// }

// func handleDecodeJSON(w http.ResponseWriter, req *http.Request) {
// 	type parameters struct {
// 		Body string `json:"body"`
// 	}

// 	decoder := json.NewDecoder(req.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	fmt.Printf("decoder: %v, params: %v", decoder, params)
// 	if err != nil {
// 		log.Printf("Error decoding parameters: %s", err)
// 		w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 		w.WriteHeader(500)
// 		return
// 	}

// 	// params is a struct with data populated successfully
// 	// ...
// }

// func handleEncodeJSON(w http.ResponseWriter, req *http.Request) {
// 	// ...

// 	type returnVals struct {
// 		Error     string    `json:"error"`
// 		CreatedAt time.Time `json:"created_at"`
// 		ID        int       `json:"id"`
// 	}

// 	respBody := returnVals{
// 		CreatedAt: time.Now(),
// 		ID:        123,
// 	}

//		dat, err := json.Marshal(respBody)
//		if err != nil {
//			log.Printf("Error marshalling JSON: %s", err)
//			w.WriteHeader(500)
//			return
//		}
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(200)
//		w.Write(dat)
//	}

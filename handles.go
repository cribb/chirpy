package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func handleOK(w http.ResponseWriter, req *http.Request) {
	text := "OK"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

func (cfg *apiConfig) addUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	email := params.Email

	user, err := cfg.dbQueries.CreateUser(req.Context(), email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(502)
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

// func respondWithSuccess(w http.ResponseWriter, code int) {
// 	type validResponse struct {
// 		Valid bool `json:"valid"`
// 	}
// 	respondWithJSON(w, code, validResponse{Valid: true})
// }

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{Error: msg})
}

func respondWithCleanVersion(w http.ResponseWriter, code int, msg string) {
	type cleanResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	cmsg := badWordsBeGone(msg)
	respondWithJSON(w, code, cleanResponse{CleanedBody: cmsg})

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
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
	cfg.dbQueries.ResetUsers(req.Context())
	text := fmt.Sprintf("Reset hits to: %v\n", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

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

// 	dat, err := json.Marshal(respBody)
// 	if err != nil {
// 		log.Printf("Error marshalling JSON: %s", err)
// 		w.WriteHeader(500)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	w.Write(dat)
// }

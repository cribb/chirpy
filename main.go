package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/cribb/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {

	var apicfg apiConfig

	// defaults to zero anyway so this isn't strictly necessary
	// but being explicit helps me right now, so here we are
	apicfg.fileserverHits.Store(0)

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error opening database.")
		return
	}

	apicfg.dbQueries = database.New(db)

	mux := http.NewServeMux()

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// fmt.Printf("Mux info: %+v\n", &mux)
	// fmt.Printf("Srv info: %+v\n", &srv)
	root := http.Dir(".")

	// var f func(http.ResponseWriter, *http.Request)

	// f = func(w http.ResponseWriter, req *http.Request) {
	// 	//fmt.Printf("w: %+v\n", w)
	// }

	fsrv := http.FileServer(root)

	mux.Handle("/app/", apicfg.middlewareMetricsInc(http.StripPrefix("/app", fsrv)))
	mux.HandleFunc("GET /api/healthz", handleOK)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("GET /admin/metrics", apicfg.getHits)
	mux.HandleFunc("POST /admin/reset", apicfg.zeroHits)

	fmt.Println("\nServer running.")
	log.Fatal(srv.ListenAndServe())

}

func handleOK(w http.ResponseWriter, req *http.Request) {
	text := "OK"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

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
			words[i] = "****"
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

func (cfg *apiConfig) zeroHits(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	text := fmt.Sprintf("Hits: %v\n", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}

// func

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

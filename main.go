package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/cribb/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	platform       string
	fileserverHits atomic.Int32
	dbURL          string
	dbQueries      *database.Queries
}

func main() {

	var apicfg apiConfig

	// defaults to zero anyway so this isn't strictly necessary
	// but being explicit helps me right now, so here we are
	apicfg.fileserverHits.Store(0)

	godotenv.Load()
	apicfg.dbURL = os.Getenv("DB_URL")
	apicfg.platform = os.Getenv("PLATFORM")

	fmt.Printf("platform: %v\n", apicfg.platform)
	fmt.Printf("Database URL: %v\n", apicfg.dbURL)

	if apicfg.platform != "dev" {
		fmt.Println("Server must be in 'dev' mode.")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", apicfg.dbURL)
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

	fsrv := http.FileServer(root)

	// Register the list of endpoints for the server/api
	mux.Handle("/app/", apicfg.middlewareMetricsInc(http.StripPrefix("/app", fsrv)))

	mux.HandleFunc("GET /api/healthz", handleOK)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("GET /admin/metrics", apicfg.getHits)
	mux.HandleFunc("POST /admin/reset", apicfg.resetAPI)
	mux.HandleFunc("POST /api/users", apicfg.addUser)
	mux.HandleFunc("POST /api/chirps", apicfg.handlePostChirp)
	mux.HandleFunc("GET /api/chirps", apicfg.handleGetAllChirps)
	mux.HandleFunc("POST /api/login", apicfg.handleUserLogin)

	mux.HandleFunc("GET /api/chirps/{chirpID}", apicfg.handleGetChirp)

	fmt.Println("\nServer running.")
	log.Fatal(srv.ListenAndServe())

}

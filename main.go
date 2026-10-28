package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	// var srv http.Server

	var apicfg apiConfig

	// defaults to zero anyway so this isn't strictly necessary
	// but being explicit helps me right now, so here we are
	apicfg.fileserverHits.Store(0)

	mux := http.NewServeMux()

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// fmt.Printf("Mux info: %+v\n", &mux)
	// fmt.Printf("Srv info: %+v\n", &srv)
	root := http.Dir(".")

	var f func(http.ResponseWriter, *http.Request)

	f = func(w http.ResponseWriter, req *http.Request) {
		text := "OK"
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Add("Scooter", "theComputer")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(text))
		//fmt.Printf("w: %+v\n", w)
	}

	fsrv := http.FileServer(root)
	// mux.Handle("/app/", http.StripPrefix("/app", fsrv))
	mux.Handle("/app/", apicfg.middlewareMetricsInc(http.StripPrefix("/app", fsrv)))
	mux.HandleFunc("GET /healthz", f)
	mux.HandleFunc("GET /metrics", apicfg.getHits)
	mux.HandleFunc("POST /reset", apicfg.zeroHits)
	// mux.Handle("/app/", middlewareLog(handler))

	fmt.Println("\nServer running.")
	log.Fatal(srv.ListenAndServe())

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
	text := fmt.Sprintf("Hits: %v\n", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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

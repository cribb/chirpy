package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// var srv http.Server

	mux := http.NewServeMux()

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Printf("Mux info: %+v\n", &mux)
	fmt.Printf("Srv info: %+v\n", &srv)
	root := http.Dir(".")

	var f func(http.ResponseWriter, *http.Request)

	f = func(w http.ResponseWriter, req *http.Request) {
		text := "OK"
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Add("Scooter", "theComputer")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(text))
		fmt.Printf("w: %+v\n", w)
	}

	fsrv := http.FileServer(root)
	mux.Handle("/app/", http.StripPrefix("/app", fsrv))
	mux.HandleFunc("/healthz", f)

	fmt.Println("\nServer running.")
	log.Fatal(srv.ListenAndServe())

}

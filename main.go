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

	// srv.Protocols = new(http.Protocols)
	// srv.Protocols.SetHTTP1(true)
	// srv.ReadTimeout = 10 * time.Second
	// srv.WriteTimeout = 10 * time.Second
	// srv.Handler = mux

	// fmt.Printf("Mux info: %+v\n", &mux)
	// fmt.Printf("Srv info: %+v\n", &srv)
	root := http.Dir(".")

	fsrv := http.FileServer(root)

	mux.Handle("/", fsrv)

	fmt.Println("\nServer running.")
	log.Fatal(srv.ListenAndServe())

}

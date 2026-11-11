package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func respondWithOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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

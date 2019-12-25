package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

var suppressLoggingPath = map[string]bool{
	"/":            true,
	"/healthcheck": true,
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !suppressLoggingPath[r.RequestURI] {
			requestDump, err := httputil.DumpRequest(r, true)
			if err != nil {
				log.Println(err)
			}
			log.Println(string(requestDump))
		}
		next.ServeHTTP(w, r)
	})
}

func responseMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(requestLogger)
	router.Use(responseMiddleWare)
	router.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	router.HandleFunc("/notes/{id}", getSingleNote).Methods("GET")
	router.HandleFunc("/notes/user/{id}", getAllNotes).Methods("GET")
	router.HandleFunc("/notes", postNote).Methods("POST")
	return router
}

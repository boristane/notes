package main

import (
	"encoding/json"
	"net/http"
)

func getSingleNote(w http.ResponseWriter, r *http.Request) {

}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still Alive")
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/mux"
)

var validate *validator.Validate

type HTTPErrorMessage struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type CreateNoteResponse struct {
	Message string `json:"message"`
	ID      uint   `json:"id"`
	URL     string `json:"url"`
}

func getSingleNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Println("Trying to access a NaN note with id:" + vars["id"])
		send500(w, "Trying to access a NaN note")
		return
	}
	note := getNote(id)
	noteJSON, err := json.Marshal(note)
	if err != nil {
		log.Println("Unable to JSON parse the note" + err.Error())
		send500(w, "")
		return
	}
	log.Printf("Got one note for id %d %s", id, string(noteJSON))
	if note.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&HTTPErrorMessage{Message: "Note not found", Code: "NONE_FOUND"})
		return
	}
	json.NewEncoder(w).Encode(&note)
}

func getAllNotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Println("Trying to access a NaN user with id:" + vars["id"])
		send500(w, "Trying to access a NaN user")
		return
	}
	notes := getNotesForUser(userID)
	notesJSON, err := json.Marshal(notes)
	if err != nil {
		log.Println("Unable to JSON parse the notes" + err.Error())
		send500(w, "")
		return
	}
	log.Printf("Got note for userID %d %s", userID, string(notesJSON))
	if len(notes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&HTTPErrorMessage{Message: "Note not found", Code: "NONE_FOUND"})
		return
	}
	json.NewEncoder(w).Encode(notes)
}

func postNote(w http.ResponseWriter, r *http.Request) {
	var requestData Note
	err := decodeAndValidateRequest(w, r, &requestData)
	if err != nil {
		return
	}
	id, createError := createNote(&requestData)
	if createError != nil {
		log.Println(createError.Error())
		send500(w, "Error creating note")
		return
	}
	json.NewEncoder(w).Encode(&CreateNoteResponse{ID: id, Message: "Note created", URL: fmt.Sprintf("/notes/%d", id)})
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still Alive")
}

func decodeAndValidateRequest(w http.ResponseWriter, r *http.Request, data interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		log.Printf("Unable to decode the request body based on the schema %v, %v. Error: %s", data, r.Body, err.Error())
	}
	validate = validator.New()
	err = validate.Struct(data)
	if err != nil {
		log.Printf("Request validation error %v. Body: %v\n", err.Error(), data)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(&HTTPErrorMessage{Message: "Error validating request", Code: "BAD_REQUEST"})
	}
	return err
}

func send500(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(&HTTPErrorMessage{Message: message, Code: "INTERNAL_SERVER_ERROR"})
}

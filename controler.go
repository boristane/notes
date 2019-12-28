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

// HTTPErrorMessage message to send when there is an error on requests
type HTTPErrorMessage struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// CreateNoteResponse the response to send to a create note request
type CreateNoteResponse struct {
	Message string `json:"message"`
	ID      uint   `json:"id"`
	URL     string `json:"url"`
}

// DeleteNoteResponse the response to send to a create note request
type DeleteNoteResponse struct {
	Message string `json:"message"`
	ID      uint64 `json:"id"`
	URL     string `json:"url"`
}

func getSingleNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		log.Println("Trying to access a NaN note with id:" + vars["id"])
		send500(w, "Trying to access a note with an invalid ID")
		return
	}
	note, errr := getNote(id)
	if errr != nil {
		log.Printf("Error getting the note, %v", err)
		send404(w)
		return
	}
	userID, ok := r.Context().Value(user_id).(uint64)
	if !ok || note.UserID != userID {
		send401(w)
		return
	}
	noteJSON, err := json.Marshal(note)
	if err != nil {
		log.Println("Unable to JSON parse the note" + err.Error())
		send500(w, "")
		return
	}
	log.Printf("Got one note for id %d %s", id, string(noteJSON))
	json.NewEncoder(w).Encode(&note)
}

func getAllNotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Println("Trying to access a NaN user with id:" + vars["id"])
		send500(w, "Trying to access a user with an invalid ID")
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
		send404(w)
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

func deleteSingleNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	userID, uerr := strconv.ParseUint(vars["userID"], 10, 32)
	if err != nil {
		log.Printf("Triying to delete a note with invalid id %s", vars["id"])
		send500(w, "Trying to delete a note with an invalid ID")
		return
	}
	if uerr != nil {
		log.Printf("Triying to delete a note with invalid userId %s", vars["id"])
		send500(w, "Trying to delete a note with an invalid ID")
		return
	}
	deleted := deleteNote(id, userID)
	if !deleted {
		log.Printf("Triying to delete a note with invalid (id, userID) pair (%d, %d)", id, userID)
		send500(w, "Unauthorised Operation")
		return
	}
	json.NewEncoder(w).Encode(&DeleteNoteResponse{ID: id, Message: "Note deleted", URL: fmt.Sprintf("/notes/user/%d", userID)})
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

func send401(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(&HTTPErrorMessage{Message: "You're not authorised to perform this action", Code: "UNAUTHORIZED"})
}

func send404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(&HTTPErrorMessage{Message: "Note not found", Code: "NONE_FOUND"})
}

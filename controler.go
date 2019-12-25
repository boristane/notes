package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type HTTPErrorMessage struct {
	Message string
	Code    string
}

func getSingleNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		panic("Trying to access a NaN note with id:" + vars["id"])
	}
	// createdNote := Note{
	// 	ID:      1,
	// 	Title:   "Je wanda22",
	// 	Content: "eNCORE PLUS",
	// 	UserID:  1,
	// }
	// createdNote2 := Note{
	// 	ID:      2,
	// 	Title:   "Je wanda",
	// 	Content: "eNCORE PLUS",
	// 	UserID:  1,
	// }
	// createNote(&createdNote)
	// createNote(&createdNote2)
	note := getNote(id)
	noteJSON, err := json.Marshal(note)
	if err != nil {
		panic("Unable to JSON parse the note" + err.Error())
	}
	log.Printf("Got one note for id %d %s", id, string(noteJSON))
	if note.UserID == 0 {
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
		panic("Trying to access a NaN user with id:" + vars["id"])
	}
	notes := getNotesForUser(userID)
	notesJSON, err := json.Marshal(notes)
	if err != nil {
		panic("Unable to JSON parse the notes" + err.Error())
	}
	log.Printf("Got note for userID %d %s", userID, string(notesJSON))
	if len(notes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(&HTTPErrorMessage{Message: "Note not found", Code: "NONE_FOUND"})
		return
	}
	json.NewEncoder(w).Encode(notes)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still Alive")
}

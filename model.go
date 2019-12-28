package main

import (
	"time"
)

// BaseModel a base model to include the on every database entity
type BaseModel struct {
	ID        uint       `json:"id" gorm:"type: serial;primary_key;column:id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`
}

// Note a representation of a user note
type Note struct {
	BaseModel
	Title    string `gorm:"not null" json:"title" validate:"required"`
	Content  string `json:"content"`
	UserID   uint64 `gorm:"not null" json:"userId" validate:"required"`
	ImageURL string `json:"imageUrl"`
}

func createNote(note *Note) (uint, error) {
	err := db.Create(&note).Error
	if err != nil {
		return 0, err
	}
	return note.ID, nil
}

func getNote(id uint64) *Note {
	var note = &Note{}
	db.First(&note, id)
	return note
}

func getNotesForUser(userID uint64) []Note {
	var notes = []Note{}
	db.Select([]string{"title", "content", "image_url"}).Where("user_id = ?", userID).Find(&notes)
	return notes
}

func updateNote(id uint, title string, content string, imageURL string) *Note {
	var note *Note
	db.First(&note, id)
	note.Content = content
	note.Title = title
	note.ImageURL = imageURL
	db.Save(&note)
	return note
}

func deleteNote(id uint64, userID uint64) bool {
	// TODO This is the worst authorisation check ever.
	var note = &Note{}
	db.First(&note, id)
	if note.UserID != userID {
		return false
	}
	db.Delete(&note)
	return true
}

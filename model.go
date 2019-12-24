package main

import (
	"github.com/jinzhu/gorm"
)

type Note struct {
	gorm.Model
	ID      uint   `gorm:"column:id;primary_index;AUTO_INCREMENT" json:"ID"`
	Title   string `gorm:"not null" json:"Title" validate:"required"`
	Content string `json:"Content"`
	UserID  int    `gorm:"not null" json:"UserId" validate:"required"`
	// CreatedAt time.Time `json:"CreatedAt" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	// UpdatedAt time.Time `json:"UpdatedAt" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	// DeletedAt time.Time `json:"DeletedAt" gorm:"default:null;column:deleted_at"`
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
	if !note.CreatedAt.IsZero() {
		note.ID = uint(id)
	}
	return note
}

func updateNote(id uint, title string, content string) *Note {
	var note *Note
	db.First(&note, id)
	note.Content = content
	note.Title = title
	db.Save(&note)
	return note
}

func deleteNote(id uint) {
	var note *Note
	db.First(&note, id)
	db.Delete(&note)
}

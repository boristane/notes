package main

import "time"

type BaseModel struct {
	ID        uint       `json:"id" gorm:"type: serial;primary_key;column:id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`
}
type Note struct {
	BaseModel
	Title   string `gorm:"not null" json:"title" validate:"required"`
	Content string `json:"content"`
	UserID  int    `gorm:"not null" json:"userId" validate:"required"`
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
	db.Where("user_id = ?", userID).Find(&notes)
	return notes
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

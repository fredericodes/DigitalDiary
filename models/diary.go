package models

import "github.com/google/uuid"

type Diary struct {
	Id      uuid.UUID
	Entries map[string]string
}

type JournalEntry struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}

func (d *Diary) ToModel(id uuid.UUID, entries map[string]string) {
	d.Id = id
	d.Entries = entries
}

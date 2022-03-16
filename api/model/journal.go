package model

import "github.com/google/uuid"

type Journal struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	EntryDate string
	Content   string
}

type JournalForm struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}

func (j *JournalForm) ToCreateModel(userId uuid.UUID) *Journal {
	return &Journal{
		Id:        uuid.New(),
		UserId:    userId,
		EntryDate: j.Date,
		Content:   j.Content,
	}
}

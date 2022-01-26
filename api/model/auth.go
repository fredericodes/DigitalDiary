package model

import "github.com/google/uuid"

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id       uuid.UUID
	Email    string
	Password string
}

func (r RegisterForm) ToModel(hashedPassword string) *User {
	return &User{
		Id:       uuid.New(),
		Email:    r.Email,
		Password: hashedPassword,
	}
}

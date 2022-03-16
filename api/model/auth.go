package model

import "github.com/google/uuid"

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id       uuid.UUID
	Username string
	Password string
}

func (r RegisterForm) ToModel(hashedPassword string) *User {
	return &User{
		Id:       uuid.New(),
		Username: r.Username,
		Password: hashedPassword,
	}
}

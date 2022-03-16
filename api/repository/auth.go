package repository

import (
	"github.com/google/uuid"

	"api/model"
)

func (db *Db) ReadUserById(id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	if err := db.Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (db *Db) ReadUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	if err := db.Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (db *Db) CreateUser(user *model.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

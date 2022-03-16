package repository

import (
	"github.com/FreddyJilson/diarynote/model"
)

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

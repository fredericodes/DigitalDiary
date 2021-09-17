package repositories

import (
	"errors"

	"github.com/go-redis/redis"
)

const (
	KeyDoesntExistsErr = "key doesnt exists"
)

func (db *Db) ReadUserExists(email string) (string, error) {
	user, err := db.Get(email).Result()
	if err != nil || err == redis.Nil {
		return "", errors.New(KeyDoesntExistsErr)
	}

	return user, nil
}

func (db *Db) CreateUser(email string, user []byte) error {
	if err := db.Set(email, user, 0).Err(); err != nil {
		return err
	}

	return nil
}

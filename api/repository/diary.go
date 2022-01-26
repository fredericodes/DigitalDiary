package repository

import (
	"github.com/google/uuid"

	"github.com/go-redis/redis"
)

func (db *Db) ReadDiaryExists(id uuid.UUID) (string, error) {
	idStr := id.String()
	diary, err := db.Get(idStr).Result()
	if err == redis.Nil {
		return "", err
	}

	return diary, nil
}

func (db *Db) CreateDiary(id uuid.UUID, user []byte) error {
	if err := db.Set(id.String(), user, 0).Err(); err != nil {
		return err
	}

	return nil
}

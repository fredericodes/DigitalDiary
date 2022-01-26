package repository

import (
	"github.com/go-redis/redis"

	"github.com/FreddyJilson/diarynote/util/configs"
)

type Db struct {
	*redis.Client
}

func New(config *configs.RedisConf) *Db {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDb,
	})

	return &Db{
		client,
	}
}

package repository

import (
	"github.com/FreddyJilson/diarynote/model"
	"github.com/FreddyJilson/diarynote/util/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	databaseConnectionErr = "failed to connect database"
)

type Db struct {
	*gorm.DB
}

func New(config *configs.DbConf) *Db {
	db, err := gorm.Open(postgres.Open(config.ToDsnString()), &gorm.Config{})
	if err != nil {
		panic(databaseConnectionErr)
	}

	return &Db{
		db,
	}
}

func (db *Db) TxBegin() *Db {
	tx := db.DB.Begin()
	return &Db{tx}
}

func (db *Db) Commit() {
	db.DB.Commit()
}

func (db *Db) Rollback() {
	db.DB.Rollback()
}

type DB interface {
	TxBegin() *Db
	Commit()
	Rollback()

	ReadUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) error
}
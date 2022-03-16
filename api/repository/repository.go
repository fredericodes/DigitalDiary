package repository

import (
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"api/model"
	"api/util/configs"
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

	ReadUserById(id uuid.UUID) (*model.User, error)
	ReadUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) error

	ReadJournalExistsByUserIdAndEntryDate(id uuid.UUID, entryDate string) (*model.Journal, error)
	CreateJournal(journal *model.Journal) error
	UpdateJournalContentByUserIdAndEntryDate(userId uuid.UUID, entryDate string, content string) (int64, error)
}

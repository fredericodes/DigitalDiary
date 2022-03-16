package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"api/model"
)

func (db *Db) ReadJournalExistsByUserIdAndEntryDate(id uuid.UUID, entryDate string) (*model.Journal, error) {
	journal := &model.Journal{}
	result := db.Table("journals j").
		Where("j.user_id = ? AND j.entry_date = ?", id, entryDate).
		Select("j.*").
		Find(journal)

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return journal, nil
}

func (db *Db) CreateJournal(journal *model.Journal) error {
	if err := db.Create(journal).Error; err != nil {
		return err
	}

	return nil
}

func (db *Db) UpdateJournalContentByUserIdAndEntryDate(userId uuid.UUID, entryDate string, content string) (int64, error) {
	result := db.Model(&model.Journal{}).
		Where("user_id = ? AND entry_date = ?", userId, entryDate).
		Updates(model.Journal{
			Content: content,
		})

	return result.RowsAffected, result.Error
}

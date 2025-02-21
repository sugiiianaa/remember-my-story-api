package repositories

import (
	"errors"

	"github.com/sugiiianaa/remember-my-story/internal/models"
	"gorm.io/gorm"
)

type JournalRepository struct {
	db *gorm.DB
}

func NewJournalRepository(db *gorm.DB) *JournalRepository {
	return &JournalRepository{db}
}

func (r *JournalRepository) Create(entry *models.JournalEntry) (uint, error) {
	if err := r.db.Create(entry).Error; err != nil {
		return 0, err
	}
	return entry.ID, nil
}

func (r *JournalRepository) FindByID(journalID uint, userID uint) (*models.JournalEntry, error) {
	var entry models.JournalEntry
	err := r.db.
		Preload("DailyTasks.SubTasks").
		Where("id = ? AND user_id = ?", journalID, userID).
		First(&entry).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("journal entry not found")
	}

	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (r *JournalRepository) UpdateEntry(journalID uint, updateData map[string]interface{}) error {
	err := r.db.Model(&models.JournalEntry{}).
		Where("id = ?", journalID).
		Updates(updateData).Error

	return err
}

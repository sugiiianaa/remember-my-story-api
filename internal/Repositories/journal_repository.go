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

func (r *JournalRepository) FindByID(id uint) (*models.JournalEntry, error) {
	var entry models.JournalEntry
	err := r.db.
		Preload("DailyTasks.SubRTasks").
		First(&entry, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("record not found")
	}

	return &entry, err
}

package repositories

import (
	"errors"
	"fmt"

	"github.com/sugiiianaa/remember-my-story/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JournalRepository struct {
	db *gorm.DB
}

func NewJournalRepository(db *gorm.DB) *JournalRepository {
	return &JournalRepository{db}
}

func (r *JournalRepository) CreateEntry(entry *models.JournalEntry) (uint, error) {
	if entry.UserID == 0 {
		return 0, errors.New("user ID must be provided")
	}

	if err := r.db.Create(entry).Error; err != nil {
		return 0, fmt.Errorf("failed to create entry: %w", err)
	}

	return entry.ID, nil
}

func (r *JournalRepository) UpdateEntry(journalID uint, userID uint, updateData map[string]interface{}) error {
	if journalID == 0 || userID == 0 {
		return errors.New("invalid input parameters")
	}

	result := r.db.Model(&models.JournalEntry{}).
		Where("id = ? AND user_id = ?", journalID, userID).
		Updates(updateData)

	if result.Error != nil {
		return fmt.Errorf("failed to update entry: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("entry not found or unauthorized")
	}

	return nil
}

func (r *JournalRepository) GetAllEntry(userID uint) ([]*models.JournalEntry, error) {
	if userID == 0 {
		return nil, errors.New("user ID must be provided")
	}

	var entries []*models.JournalEntry
	err := r.db.
		Preload("DailyTasks.SubTasks").
		Where("user_id = ?", userID).
		Find(&entries).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve entries: %w", err)
	}

	return entries, nil
}

func (r *JournalRepository) FindEntryByID(journalID uint, userID uint) (*models.JournalEntry, error) {
	if journalID == 0 || userID == 0 {
		return nil, errors.New("invalid input parameters")
	}

	var entry models.JournalEntry
	err := r.db.
		Preload("DailyTasks.SubTasks").
		Where("id = ? AND user_id = ?", journalID, userID).
		First(&entry).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("entry not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve entry: %w", err)
	}

	return &entry, nil
}

func (r *JournalRepository) DeleteEntry(journalID, userID uint) error {
	if journalID == 0 || userID == 0 {
		return errors.New("invalid input parameters")
	}

	// Explicitly load nested relationships
	var entry models.JournalEntry
	err := r.db.
		Preload("DailyTasks.SubTasks").
		Where("id = ? AND user_id = ?", journalID, userID).
		First(&entry).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("entry not found or unauthorized")
		}
		return fmt.Errorf("failed to retrieve entry: %w", err)
	}

	// Delete with all associations
	result := r.db.Select(clause.Associations).Delete(&entry)

	if result.Error != nil {
		return fmt.Errorf("failed to delete entry: %w", result.Error)
	}

	return nil
}

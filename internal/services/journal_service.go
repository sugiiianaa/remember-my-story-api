package services

import (
	"time"

	repositories "github.com/sugiiianaa/remember-my-story/internal/Repositories"
	"github.com/sugiiianaa/remember-my-story/internal/models"
)

type JournalService struct {
	journalRepo *repositories.JournalRepository
}

func NewJournalService(journalRepo *repositories.JournalRepository) *JournalService {
	return &JournalService{journalRepo: journalRepo}
}

func (s *JournalService) CreateEntry(entry *models.JournalEntry) (uint, error) {
	// Set the date to the beginning of the day
	entry.Date = time.Date(entry.Date.Year(), entry.Date.Month(), entry.Date.Day(),
		0, 0, 0, 0, entry.Date.Location())

	return s.journalRepo.Create(entry)
}

func (s *JournalService) UpdateEntry(journalID uint, userID uint, updateData map[string]interface{}) error {
	// Check if entry exists and belongs to the user
	existingEntry, err := s.journalRepo.FindByID(journalID, userID)
	if err != nil {
		return err
	}

	// Update the entry using the repository
	err = s.journalRepo.UpdateEntry(existingEntry.ID, updateData)
	if err != nil {
		return err
	}

	return nil
}

package services

import (
	"context"
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

func (s *JournalService) GetEntry(ctx context.Context, id uint) (*models.JournalEntry, error) {
	return s.journalRepo.FindByID(id)
}

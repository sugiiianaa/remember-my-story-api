package services

import (
	"context"
	"time"

	repositories "github.com/sugiiianaa/remember-my-story/internal/Repositories"
	"github.com/sugiiianaa/remember-my-story/internal/models"
)

type JournalService interface {
	CreateEntry(ctx context.Context, entry *models.JournalEntry) error
	GetEntry(ctx context.Context, id uint) (*models.JournalEntry, error)
	UpdateEntry(ctx context.Context, entry *models.JournalEntry) error
	GetEntriesByDate(ctx context.Context, date time.Time) ([]models.JournalEntry, error)
	GetAllEntries(ctx context.Context, userID uint) ([]models.JournalEntry, error)
}

type journalService struct {
	repo repositories.JournalRepository
}

func NewJournalService(repo repositories.JournalRepository) JournalService {
	return &journalService{repo: repo}
}

func (s *journalService) CreateEntry(ctx context.Context, entry *models.JournalEntry) error {
	// Set the date to the beginning of the day
	entry.Date = time.Date(entry.Date.Year(), entry.Date.Month(), entry.Date.Day(),
		0, 0, 0, 0, entry.Date.Location())

	// Basic validation
	if entry.ThisDayDescription == "" {
		return ValidationError{"this_day_description is required"}
	}

	return s.repo.Create(ctx, entry)
}

func (s *journalService) GetEntry(ctx context.Context, id uint) (*models.JournalEntry, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *journalService) UpdateEntry(ctx context.Context, entry *models.JournalEntry) error {
	existingEntry, err := s.repo.FindByID(ctx, entry.ID)
	if err != nil {
		return err
	}

	// Prevent date modification
	entry.Date = existingEntry.Date
	return s.repo.Update(ctx, entry)
}

func (s *journalService) GetEntriesByDate(ctx context.Context, date time.Time) ([]models.JournalEntry, error) {
	return s.repo.FindByDate(ctx, date)
}

func (s *journalService) GetAllEntries(ctx context.Context, userID uint) ([]models.JournalEntry, error) {
	return s.repo.FindAll(ctx, userID)
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

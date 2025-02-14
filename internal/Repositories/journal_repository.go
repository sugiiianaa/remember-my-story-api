package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/sugiiianaa/remember-my-story/internal/models"
	"gorm.io/gorm"
)

type JournalRepository interface {
	Create(ctx context.Context, entry *models.JournalEntry) error
	FindByID(ctx context.Context, id uint) (*models.JournalEntry, error)
	Update(ctx context.Context, entry *models.JournalEntry) error
	FindByDate(ctx context.Context, date time.Time) ([]models.JournalEntry, error)
}

type journalRepository struct {
	db *gorm.DB
}

func NewJournalRepository(db *gorm.DB) JournalRepository {
	return &journalRepository{db: db}
}

func (r *journalRepository) Create(ctx context.Context, entry *models.JournalEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *journalRepository) FindByID(ctx context.Context, id uint) (*models.JournalEntry, error) {
	var entry models.JournalEntry
	err := r.db.
		WithContext(ctx).
		Preload("DailyTasks.SubRTasks").
		First(&entry, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return &entry, err
}

func (r *journalRepository) Update(ctx context.Context, entry *models.JournalEntry) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(entry).Error
}

func (r *journalRepository) FindByDate(ctx context.Context, date time.Time) ([]models.JournalEntry, error) {
	var entries []models.JournalEntry
	err := r.db.WithContext(ctx).
		Where("date = ?", date.Format("2006-01-02")).
		Find(&entries).Error
	return entries, err
}

var ErrNotFound = errors.New("record not found")

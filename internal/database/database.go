package database

import (
	"fmt"

	"github.com/sugiiianaa/remember-my-story/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(host, user, password, dbname, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&models.JournalEntry{},
		&models.DailyTask{},
		&models.DailySubTask{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	// Add foreign key constraints if they don't exist
	if !db.Migrator().HasConstraint(&models.DailyTask{}, "JournalEntry") {
		err = db.Migrator().CreateConstraint(&models.DailyTask{}, "JournalEntry")
		if err != nil {
			return nil, fmt.Errorf("failed to create JournalEntry constraint: %w", err)
		}
	}

	if !db.Migrator().HasConstraint(&models.DailySubTask{}, "DailyTask") {
		err = db.Migrator().CreateConstraint(&models.DailySubTask{}, "DailyTask")
		if err != nil {
			return nil, fmt.Errorf("failed to create DailyTask constraint: %w", err)
		}
	}

	return db, nil
}

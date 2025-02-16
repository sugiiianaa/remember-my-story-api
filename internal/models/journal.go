package models

import (
	"time"

	"github.com/sugiiianaa/remember-my-story/internal/models/enums"
	"gorm.io/gorm"
)

type JournalEntry struct {
	gorm.Model
	Date               time.Time      `gorm:"not null; index"`
	Mood               enums.MoodType `gorm:"not null; index"`
	ThisDayDescription string         `gorm:"not null"`
	DailyReflection    string         `gorm:"not null"`
	UserID             uint           `gorm:"not null; index"`
	DailyTasks         []DailyTask    `gorm:"foreignKey:JournalEntryID"`
}

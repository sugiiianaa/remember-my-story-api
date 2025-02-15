package models

import (
	"time"

	"gorm.io/gorm"
)

type JournalEntry struct {
	gorm.Model
	Date               time.Time `gorm:"index"`
	Mood               int
	ThisDayDescription string
	DailyReflection    string
	UserID             uint        `gorm:"foreignKey:UserID"`
	DailyTasks         []DailyTask `gorm:"foreignKey:JournalEntryID"`
}

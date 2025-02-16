package models

import (
	"time"

	"github.com/sugiiianaa/remember-my-story/internal/models/enums"
	"gorm.io/gorm"
)

type JournalEntry struct {
	gorm.Model
	Date               time.Time `gorm:"index"`
	Mood               enums.MoodType
	ThisDayDescription string
	DailyReflection    string
	UserID             uint        `gorm:"foreignKey:UserID" json:"-"`
	DailyTasks         []DailyTask `gorm:"foreignKey:JournalEntryID"`
}

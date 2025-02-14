package models

import (
	"gorm.io/gorm"
)

type DailyTask struct {
	gorm.Model
	JournalEntryID uint `gorm:"index"` // Add index for better performance
	Task           string
	Status         bool
	SubTasks       []DailySubTask `gorm:"foreignKey:DailyTaskID"`
}

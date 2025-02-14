package models

import "gorm.io/gorm"

type DailySubTask struct {
	gorm.Model
	DailyTaskID uint `gorm:"index"` // Changed from TaskId to DailyTaskID
	SubTask     string
	Status      bool
}

package repositories

import (
	"errors"
	"fmt"

	"github.com/sugiiianaa/remember-my-story/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO : Add CRUD daily task repo
type DailyTaskRepository struct {
	db *gorm.DB
}

func NewDailyTaskRepostiory(db *gorm.DB) *DailyTaskRepository {
	return &DailyTaskRepository{db}
}

func (r *DailyTaskRepository) CreateDailyTask(userID uint, dailyTask *models.DailyTask) (uint, error) {
	if userID == 0 {
		return 0, errors.New("user ID must be provided")
	}

	if err := r.db.Create(dailyTask).Error; err != nil {
		return 0, fmt.Errorf("failed to create daily task: %w", err)
	}
	return dailyTask.ID, nil
}

func (r *DailyTaskRepository) UpdateDailyTask(dailyTaskID uint, userID uint, updateData map[string]interface{}) error {
	if dailyTaskID == 0 || userID == 0 {
		return errors.New("invalid input parameters")
	}

	result := r.db.Model(&models.DailyTask{}).
		Where("id = ? AND user_id = ?", dailyTaskID, userID).
		Updates(updateData)

	if result.Error != nil {
		return fmt.Errorf("failed to update daily task: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("daily task not found or unauthorized")
	}

	return nil
}

func (r *DailyTaskRepository) FindDailyTaskByID(dailyTaskID uint, userID uint) (*models.DailyTask, error) {
	if dailyTaskID == 0 || userID == 0 {
		return nil, errors.New("invalid input parameters")
	}

	var dailyTask models.DailyTask
	err := r.db.
		Preload("SubTasks").
		Where("id = ? AND user_id = ?", dailyTaskID, userID).
		First(&dailyTask).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("entry not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve daily task: %w", err)
	}

	return &dailyTask, nil
}

func (r *DailyTaskRepository) DeleteDailyTask(dailyTaskID uint, userID uint) error {
	if dailyTaskID == 0 || userID == 0 {
		return errors.New("invalid input parameters")
	}

	// Explicitly load nested relationships
	var entry models.DailyTask
	err := r.db.
		Preload("SubTasks").
		Where("id = ? AND user_id = ?", dailyTaskID, userID).
		First(&entry).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("edaily task found or unauthorized")
		}
		return fmt.Errorf("failed to retrieve daily task: %w", err)
	}

	// Delete with all associations
	result := r.db.Select(clause.Associations).Delete(&entry)

	if result.Error != nil {
		return fmt.Errorf("failed to daily task : %w", result.Error)
	}

	return nil
}

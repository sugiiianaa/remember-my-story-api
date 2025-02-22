package services

import (
	repositories "github.com/sugiiianaa/remember-my-story/internal/Repositories"
	"github.com/sugiiianaa/remember-my-story/internal/models"
)

// TODO : Add CRUD daily task service
type DailyTaskService struct {
	dailyTaskRepo *repositories.DailyTaskRepository
}

func NewDailyTaskService(dailyTaskRepo *repositories.DailyTaskRepository) *DailyTaskService {
	return &DailyTaskService{dailyTaskRepo: dailyTaskRepo}
}

func (s *DailyTaskService) AddDailyTask(userID uint, dailyTask *models.DailyTask) (uint, error) {
	dailyTask.Status = false
	return s.dailyTaskRepo.CreateDailyTask(userID, dailyTask)
}

func (s *DailyTaskService) UpdateDailyTask(dailyTaskID uint, userID uint, updateData map[string]interface{}) error {
	// Check if entry exists and belongs to the user
	existingDailyTask, err := s.dailyTaskRepo.FindDailyTaskByID(dailyTaskID, userID)
	if err != nil {
		return err
	}

	// Update the entry using the repository
	err = s.dailyTaskRepo.UpdateDailyTask(existingDailyTask.ID, userID, updateData)
	if err != nil {
		return err
	}

	return nil
}

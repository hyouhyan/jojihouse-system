package service

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
)

type RemainingEntriesLogService struct {
	repo *repository.RemainingEntriesLogRepository
}

func NewRemainingEntriesLogService(repo *repository.RemainingEntriesLogRepository) *RemainingEntriesLogService {
	return &RemainingEntriesLogService{repo: repo}
}

func (s *RemainingEntriesLogService) CreateRemainingEntriesLog(log *model.RemainingEntriesLog) error {
	return s.repo.CreateRemainingEntriesLog(log)
}

func (s *RemainingEntriesLogService) GetRemainingEntriesLogs() ([]model.RemainingEntriesLog, error) {
	return s.repo.GetRemainingEntriesLogs()
}

func (s *RemainingEntriesLogService) GetRemainingEntriesLogsByUserID(userID int) ([]model.RemainingEntriesLog, error) {
	return s.repo.GetRemainingEntriesLogsByUserID(userID)
}

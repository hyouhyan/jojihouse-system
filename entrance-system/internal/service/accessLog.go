package service

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
	"time"
)

type LogService struct {
	repo *repository.LogRepository
}

func NewLogService(repo *repository.LogRepository) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) CreateAccessLog(log *model.AccessLog) error {
	return s.repo.CreateAccessLog(log)
}

func (s *LogService) CreateEntryAccessLog(userid int) error {
	log :=
		&model.AccessLog{
			UserID:     userid,
			Time:       time.Now(),
			AccessType: "entry",
		}

	return s.CreateAccessLog(log)
}

func (s *LogService) CreateExitAccessLog(userid int) error {
	log :=
		&model.AccessLog{
			UserID:     userid,
			Time:       time.Now(),
			AccessType: "exit",
		}

	return s.CreateAccessLog(log)
}

func (s *LogService) GetAccessLogs() ([]model.AccessLog, error) {
	return s.repo.GetAccessLogs()
}

func (s *LogService) GetAccessLogsByUserID(userID int) ([]model.AccessLog, error) {
	return s.repo.GetAccessLogsByUserID(userID)
}

func (s *LogService) CreateRemainingEntriesLog(log *model.RemainingEntriesLog) error {
	return s.repo.CreateRemainingEntriesLog(log)
}

func (s *LogService) GetRemainingEntriesLogs() ([]model.RemainingEntriesLog, error) {
	return s.repo.GetRemainingEntriesLogs()
}

func (s *LogService) GetRemainingEntriesLogsByUserID(userID int) ([]model.RemainingEntriesLog, error) {
	return s.repo.GetRemainingEntriesLogsByUserID(userID)
}

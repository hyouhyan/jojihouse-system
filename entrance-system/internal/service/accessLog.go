package service

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
	"time"
)

type AccessLogService struct {
	repo *repository.AccessLogRepository
}

func NewAccessLogService(repo *repository.AccessLogRepository) *AccessLogService {
	return &AccessLogService{repo: repo}
}

func (s *AccessLogService) CreateAccessLog(log *model.AccessLog) error {
	return s.repo.CreateAccessLog(log)
}

func (s *AccessLogService) CreateEntryAccessLog(userid int) error {
	log :=
		&model.AccessLog{
			UserID:     userid,
			Time:       time.Now(),
			AccessType: "entry",
		}

	return s.CreateAccessLog(log)
}

func (s *AccessLogService) CreateExitAccessLog(userid int) error {
	log :=
		&model.AccessLog{
			UserID:     userid,
			Time:       time.Now(),
			AccessType: "exit",
		}

	return s.CreateAccessLog(log)
}

func (s *AccessLogService) GetAccessLogs() ([]model.AccessLog, error) {
	return s.repo.GetAccessLogs()
}

func (s *AccessLogService) GetAccessLogsByUserID(userID int) ([]model.AccessLog, error) {
	return s.repo.GetAccessLogsByUserID(userID)
}

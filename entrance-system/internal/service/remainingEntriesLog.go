package service

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RemainingEntriesLogService struct {
	repo *repository.RemainingEntriesLogRepository
}

func NewRemainingEntriesLogService(repo *repository.RemainingEntriesLogRepository) *RemainingEntriesLogService {
	return &RemainingEntriesLogService{repo: repo}
}

func (s *RemainingEntriesLogService) CreateRemainingEntriesLog(log *model.RemainingEntriesLog) error {
	log.ID = primitive.NilObjectID
	log.UpdatedAt = time.Now()
	return s.repo.CreateRemainingEntriesLog(log)
}

func (s *RemainingEntriesLogService) GetRemainingEntriesLogs(lastID primitive.ObjectID) ([]model.RemainingEntriesLog, error) {
	return s.repo.GetRemainingEntriesLogs(lastID)
}

func (s *RemainingEntriesLogService) GetRemainingEntriesLogsByUserID(userID int, lastID primitive.ObjectID) ([]model.RemainingEntriesLog, error) {
	return s.repo.GetRemainingEntriesLogsByUserID(userID, lastID)
}

func (s *RemainingEntriesLogService) GetLastRemainingEntriesLogByUserID(userID int) (model.RemainingEntriesLog, error) {
	return s.repo.GetLastRemainingEntriesLogByUserID(userID)
}

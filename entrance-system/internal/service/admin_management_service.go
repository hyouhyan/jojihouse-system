package service

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminManagementService struct {
	userRepository                *repository.UserRepository
	roleRepository                *repository.RoleRepository
	accessLogRepository           *repository.AccessLogRepository
	remainingEntriesLogRepository *repository.RemainingEntriesLogRepository
}

func NewAdminManagementService(userRepository *repository.UserRepository, roleRepository *repository.RoleRepository, accessLogRepository *repository.AccessLogRepository, remainingEntriesLogRepository *repository.RemainingEntriesLogRepository) *AdminManagementService {
	return &AdminManagementService{userRepository: userRepository, roleRepository: roleRepository, accessLogRepository: accessLogRepository, remainingEntriesLogRepository: remainingEntriesLogRepository}
}

func (s *AdminManagementService) CreateUser(user *model.User) (*model.User, error) {
	return s.userRepository.CreateUser(user)
}

func (s *AdminManagementService) UpdateUser(user *model.User) error {
	return s.userRepository.UpdateUser(user)
}

func (s *AdminManagementService) DeleteUser(userID int) error {
	return s.userRepository.DeleteUser(userID)
}

func (s *AdminManagementService) IncreaseRemainingEntries(userID int, count int) error {
	return s.userRepository.IncreaseRemainingEntries(userID, count)
}

func (s *AdminManagementService) AddRoleToUser(userID, roleID int) error {
	return s.roleRepository.AddRoleToUser(userID, roleID)
}

func (s *AdminManagementService) RemoveRoleFromUser(userID, roleID int) error {
	return s.roleRepository.RemoveRoleFromUser(userID, roleID)
}

func (s *AdminManagementService) GetAccessLogs(lastID primitive.ObjectID) ([]model.AccessLog, error) {
	return s.accessLogRepository.GetAccessLogs(lastID)
}

func (s *AdminManagementService) GetRemainingEntriesLogs(lastID primitive.ObjectID) ([]model.RemainingEntriesLog, error) {
	return s.remainingEntriesLogRepository.GetRemainingEntriesLogs(lastID)
}

func (s *AdminManagementService) GetRemainingEntriesLogsOnlyIncrease(lastID primitive.ObjectID) ([]model.RemainingEntriesLog, error) {
	return s.remainingEntriesLogRepository.GetRemainingEntriesLogsOnlyIncrease(lastID)
}
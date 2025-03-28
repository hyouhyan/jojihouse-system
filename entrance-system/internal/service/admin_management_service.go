package service

import (
	"jojihouse-entrance-system/internal/repository"
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

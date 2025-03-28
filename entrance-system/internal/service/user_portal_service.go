package service

import (
	"jojihouse-entrance-system/internal/repository"
)

type UserPortalService struct {
	userRepository                *repository.UserRepository
	roleRepository                *repository.RoleRepository
	accessLogRepository           *repository.AccessLogRepository
	remainingEntriesLogRepository *repository.RemainingEntriesLogRepository
}

func NewUserPortalService(userRepository *repository.UserRepository, roleRepository *repository.RoleRepository, accessLogRepository *repository.AccessLogRepository, remainingEntriesLogRepository *repository.RemainingEntriesLogRepository) *UserPortalService {
	return &UserPortalService{userRepository: userRepository, roleRepository: roleRepository, accessLogRepository: accessLogRepository, remainingEntriesLogRepository: remainingEntriesLogRepository}
}

package service

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) GetRoleByID(id int) (*model.Role, error) {
	return s.roleRepo.GetRoleByID(id)
}

func (s *RoleService) GetRoleByName(name string) (*model.Role, error) {
	return s.roleRepo.GetRoleByName(name)
}

func (s *RoleService) GetRolesByUserID(userID int) ([]model.Role, error) {
	return s.roleRepo.GetRolesByUserID(userID)
}

func (s *RoleService) IsStudent(userID int) (bool, error) {
	return s.roleRepo.IsStudent(userID)
}

func (s *RoleService) IsHouseAdmin(userID int) (bool, error) {
	return s.roleRepo.IsHouseAdmin(userID)
}

func (s *RoleService) IsSystemAdmin(userID int) (bool, error) {
	return s.roleRepo.IsSystemAdmin(userID)
}

func (s *RoleService) IsGuest(userID int) (bool, error) {
	return s.roleRepo.IsGuest(userID)
}

func (s *RoleService) AddRoleToUser(userID int, roleID int) error {
	return s.roleRepo.AddRoleToUser(userID, roleID)
}

func (s *RoleService) RemoveRoleFromUser(userID int, roleID int) error {
	return s.roleRepo.RemoveRoleFromUser(userID, roleID)
}

package service

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
	"time"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(id int) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) GetUserByBarcode(barcode string) (*model.User, error) {
	return s.repo.GetUserByBarcode(barcode)
}

func (s *UserService) CreateUser(user *model.User) (*model.User, error) {
	user.Registered_at = time.Now()
	return s.repo.CreateUser(user)
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

func (s *UserService) DecreaseRemainingEntries(id int) error {
	return s.repo.DecreaseRemainingEntries(id)
}

func (s *UserService) IncreaseRemainingEntries(id int, count int) error {
	return s.repo.IncreaseRemainingEntries(id, count)
}

func (s *UserService) IncreaseTotalEntries(id int) error {
	return s.repo.IncreaseTotalEntries(id)
}

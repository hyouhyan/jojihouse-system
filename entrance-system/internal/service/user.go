package service

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
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
	return s.repo.CreateUser(user)
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

// 入場したときの処理
func (s *UserService) EnterUser(barcode string) error {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.repo.GetUserByBarcode(barcode)
	if err != nil {
		return err
	}

	// TODO: ログの生成

	// TODO: ユーザーが入場可能回数を減らす対象かの確認
	// ハウス管理者とか、同日の再入場とか

	// 残り回数を減らす
	err = s.repo.DecreaseRemainingEntries(user.ID)
	if err != nil {
		return err
	}

	return nil
}

// 退場したときの処理
func (s *UserService) ExitUser(barcode string) error {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.repo.GetUserByBarcode(barcode)
	if err != nil {
		return err
	}

	// TODO: ログの生成

	// 入場回数を増やす
	err = s.repo.IncreaseTotalEntries(user.ID)
	if err != nil {
		return err
	}

	return nil
}

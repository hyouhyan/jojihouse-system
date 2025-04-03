package service

import (
	"jojihouse-entrance-system/api/model/request"
	"jojihouse-entrance-system/api/model/response"
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
	"time"

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

func (s *AdminManagementService) CreateUser(req *request.CreateUserRequest) (*response.UserResponse, error) {
	// パース的な、model.userに合わせて再構築
	user := &model.User{
		Name:              req.Name,
		Description:       req.Description,
		Barcode:           req.Barcode,
		Contact:           req.Contact,
		Remaining_entries: req.Remaining_entries,
	}

	// ユーザーを作成
	user, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	res := &response.UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		Description:       user.Description,
		Barcode:           user.Barcode,
		Contact:           user.Contact,
		Remaining_entries: user.Remaining_entries,
		Registered_at:     user.Registered_at,
		Total_entries:     user.Total_entries,
	}

	return res, nil
}

func (s *AdminManagementService) UpdateUser(userID int, user *request.UpdateUser) error {
	userModel := &model.User{
		ID:                userID,
		Name:              user.Name,
		Description:       user.Description,
		Barcode:           user.Barcode,
		Contact:           user.Contact,
		Remaining_entries: user.Remaining_entries,
	}

	return s.userRepository.UpdateUser(userModel)
}

func (s *AdminManagementService) DeleteUser(userID int) error {
	return s.userRepository.DeleteUser(userID)
}

// 入場可能回数の追加
func (s *AdminManagementService) IncreaseRemainingEntries(userID int, count int, reason string, updatedBy string) error {
	// 入場可能回数 追加
	beforeCount, afterCount, err := s.userRepository.IncreaseRemainingEntries(userID, count)
	if err != nil {
		return err
	}

	// ログに保存
	log := &model.RemainingEntriesLog{
		UserID:          userID,
		PreviousEntries: beforeCount,
		NewEntries:      afterCount,
		Reason:          reason,
		UpdatedBy:       updatedBy,
		UpdatedAt:       time.Now(),
	}

	return s.remainingEntriesLogRepository.CreateRemainingEntriesLog(log)
}

func (s *AdminManagementService) AddRoleToUser(userID, roleID int) error {
	return s.roleRepository.AddRoleToUser(userID, roleID)
}

func (s *AdminManagementService) RemoveRoleFromUser(userID, roleID int) error {
	return s.roleRepository.RemoveRoleFromUser(userID, roleID)
}

func (s *AdminManagementService) GetAllAccessLogs(lastID primitive.ObjectID) ([]model.AccessLog, error) {
	// フィルターをあえて指定しない
	options := model.AccessLogFilter{}

	return s.accessLogRepository.GetAccessLogsByAnyFilter(lastID, options)
}

func (s *AdminManagementService) GetRemainingEntriesLogs(lastID primitive.ObjectID) ([]model.RemainingEntriesLog, error) {
	return s.remainingEntriesLogRepository.GetRemainingEntriesLogs(lastID, 50)
}

func (s *AdminManagementService) GetRemainingEntriesLogsOnlyIncrease(lastID primitive.ObjectID) ([]model.RemainingEntriesLog, error) {
	return s.remainingEntriesLogRepository.GetRemainingEntriesLogsOnlyIncrease(lastID, 50)
}

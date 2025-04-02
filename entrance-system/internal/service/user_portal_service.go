package service

import (
	"errors"
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
	"jojihouse-entrance-system/internal/response"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPortalService struct {
	userRepository                *repository.UserRepository
	roleRepository                *repository.RoleRepository
	accessLogRepository           *repository.AccessLogRepository
	remainingEntriesLogRepository *repository.RemainingEntriesLogRepository
	currentUsersRepository        *repository.CurrentUsersRepository
}

func NewUserPortalService(userRepository *repository.UserRepository,
	roleRepository *repository.RoleRepository,
	accessLogRepository *repository.AccessLogRepository,
	remainingEntriesLogRepository *repository.RemainingEntriesLogRepository,
	currentUsersRepository *repository.CurrentUsersRepository,
) *UserPortalService {
	return &UserPortalService{userRepository: userRepository,
		roleRepository:                roleRepository,
		accessLogRepository:           accessLogRepository,
		remainingEntriesLogRepository: remainingEntriesLogRepository,
		currentUsersRepository:        currentUsersRepository,
	}
}

// ログの取得
func (s *UserPortalService) GetAccessLogsByUserID(userID int, lastID primitive.ObjectID) ([]response.AccessLogResponse, error) {
	options := model.AccessLogFilter{
		UserID: userID,
	}
	return s.GetAccessLogsByAnyFilter(lastID, options)
}

func (s *UserPortalService) GetAccessLogsByAnyFilter(lastID primitive.ObjectID, options ...model.AccessLogFilter) ([]response.AccessLogResponse, error) {
	opt := model.AccessLogFilter{}

	if len(options) > 0 {
		opt = options[0]
		// Limitの上限を50に
		if opt.Limit > 50 || opt.Limit <= 0 {
			opt.Limit = 50
		}

		// UserIDが正しいか
		if opt.UserID < 0 {
			return nil, errors.New("UserIDが正しくありません")
		}

		// DayBeforeとDayAfterの整合性
		if !opt.DayBefore.IsZero() && !opt.DayAfter.IsZero() && opt.DayBefore.After(opt.DayAfter) {
			return nil, errors.New("DayBefore cannot be after DayAfter")
		}
	}

	logs, err := s.accessLogRepository.GetAccessLogsByAnyFilter(lastID, opt)
	if err != nil {
		return nil, err
	}

	// UserIDの一覧を作成
	userIDs := make([]int, len(logs))
	for i, log := range logs {
		userIDs[i] = log.UserID
	}

	// PostgreSQL から UserID に対応する UserName を取得
	users, err := s.userRepository.GetUsersByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	// UserID -> UserName のマッピング
	userMap := make(map[int]string)
	for _, user := range users {
		userMap[user.ID] = user.Name
	}

	// レスポンスデータを作成
	var responseLogs []response.AccessLogResponse
	for _, log := range logs {
		responseLogs = append(responseLogs, response.AccessLogResponse{
			ID:         log.ID.Hex(),
			UserID:     log.UserID,
			UserName:   userMap[log.UserID], // UserIDからUserNameを取得
			Time:       log.Time,
			AccessType: log.AccessType,
		})
	}

	return responseLogs, nil
}

func (s *UserPortalService) GetRemainingEntriesLogsByUserID(userID int, lastID primitive.ObjectID) ([]model.RemainingEntriesLog, error) {
	return s.remainingEntriesLogRepository.GetRemainingEntriesLogsByUserID(userID, lastID, 50)
}

// ユーザー情報の取得
func (s *UserPortalService) GetUserByID(userID int) (*model.User, error) {
	return s.userRepository.GetUserByID(userID)
}

func (s *UserPortalService) GetUserByBarcode(barcode string) (*model.User, error) {
	return s.userRepository.GetUserByBarcode(barcode)
}

// ロール関連
func (s *UserPortalService) GetRolesByUserID(userID int) ([]model.Role, error) {
	return s.roleRepository.GetRolesByUserID(userID)
}

func (s *UserPortalService) GetRoleByID(roleID int) (*model.Role, error) {
	return s.roleRepository.GetRoleByID(roleID)
}

func (s *UserPortalService) GetRoleByName(name string) (*model.Role, error) {
	return s.roleRepository.GetRoleByName(name)
}

func (s *UserPortalService) IsMember(userID int) (bool, error) {
	return s.roleRepository.IsMember(userID)
}

func (s *UserPortalService) IsStudent(userID int) (bool, error) {
	return s.roleRepository.IsStudent(userID)
}

func (s *UserPortalService) IsHouseAdmin(userID int) (bool, error) {
	return s.roleRepository.IsHouseAdmin(userID)
}

func (s *UserPortalService) IsSystemAdmin(userID int) (bool, error) {
	return s.roleRepository.IsSystemAdmin(userID)
}

func (s *UserPortalService) IsGuest(userID int) (bool, error) {
	return s.roleRepository.IsGuest(userID)
}

// 在室ユーザー一覧を取得
func (s *UserPortalService) GetCurrentUsers() ([]model.CurrentUser, error) {
	return s.currentUsersRepository.GetCurrentUsers()
}

// 入室時間を取得
func (s *UserPortalService) GetEnteredTime(userID int) (time.Time, error) {
	return s.currentUsersRepository.GetEnteredTime(userID)
}

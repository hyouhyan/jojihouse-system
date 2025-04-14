package service

import (
	"errors"
	"jojihouse-entrance-system/api/model/response"
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
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
func (s *UserPortalService) GetAccessLogsByUserID(userID int, lastID string) ([]response.AccessLog, error) {
	options := model.AccessLogFilter{
		UserID: userID,
	}
	return s.GetAccessLogsByAnyFilter(lastID, options)
}

func (s *UserPortalService) GetAccessLogsByAnyFilter(lastID string, options ...model.AccessLogFilter) ([]response.AccessLog, error) {
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
		if !opt.DayBefore.IsZero() && !opt.DayAfter.IsZero() {
			if opt.DayBefore.After(opt.DayAfter) {
				return nil, errors.New("DayBefore cannot be after DayAfter")
			}

			// DayAfterの時間を00:00:00に設定
			opt.DayAfter = time.Date(opt.DayAfter.Year(), opt.DayAfter.Month(), opt.DayAfter.Day(), 0, 0, 0, 0, opt.DayAfter.Location())

			// DayBeforeの時間を翌日の00:00:00に設定
			opt.DayBefore = time.Date(opt.DayBefore.Year(), opt.DayBefore.Month(), opt.DayBefore.Day()+1, 0, 0, 0, 0, opt.DayBefore.Location())
		}
	}

	var objectID primitive.ObjectID
	var err error

	// lastIDを変換
	if lastID == "" {
		objectID = primitive.NilObjectID
	} else {
		objectID, err = primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
	}

	logs, err := s.accessLogRepository.GetAccessLogsByAnyFilter(objectID, opt)
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
	var responseLogs []response.AccessLog
	for _, log := range logs {
		responseLogs = append(responseLogs, response.AccessLog{
			ID:         log.ID.Hex(),
			UserID:     log.UserID,
			UserName:   userMap[log.UserID], // UserIDからUserNameを取得
			Time:       log.Time,
			AccessType: log.AccessType,
		})
	}

	return responseLogs, nil
}

func (s *UserPortalService) GetRemainingEntriesLogsByUserID(userID int, lastID string) ([]response.RemainingEntriesLog, error) {
	var objectID primitive.ObjectID
	var err error
	// lastIDを変換
	if lastID == "" {
		objectID = primitive.NilObjectID
	} else {
		objectID, err = primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
	}

	logs, err := s.remainingEntriesLogRepository.GetRemainingEntriesLogsByUserID(userID, objectID, 50)
	if err != nil {
		return nil, err
	}

	// レスポンスデータの作成
	var responseLogs []response.RemainingEntriesLog
	for _, log := range logs {
		responseLogs = append(responseLogs, response.RemainingEntriesLog{
			ID:              log.ID.Hex(),
			UserID:          log.UserID,
			PreviousEntries: log.PreviousEntries,
			NewEntries:      log.NewEntries,
			Reason:          log.Reason,
			UpdatedBy:       log.UpdatedBy,
			UpdatedAt:       log.UpdatedAt,
		})
	}
	return responseLogs, nil
}

func (s *UserPortalService) GetAllUsers() ([]response.User, error) {
	users, err := s.userRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var res []response.User
	for _, user := range users {
		res = append(res, *s.cnvModelUserToResponseUser(&user))
	}

	return res, err
}

// ユーザー情報の取得
func (s *UserPortalService) GetUserByID(userID int) (*response.User, error) {
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.cnvModelUserToResponseUser(user), nil
}

func (s *UserPortalService) GetUserByBarcode(barcode string) (*model.User, error) {
	return s.userRepository.GetUserByBarcode(barcode)
}

func (s *UserPortalService) GetAllRoles() ([]response.Role, error) {
	roles, err := s.roleRepository.GetAllRoles()
	if err != nil {
		return nil, err
	}

	var res []response.Role
	for _, role := range roles {
		res = append(res, response.Role{
			ID:   role.ID,
			Name: role.Name,
		})
	}
	return res, nil
}

// ロール関連
func (s *UserPortalService) GetRolesByUserID(userID int) ([]response.Role, error) {
	roles, err := s.roleRepository.GetRolesByUserID(userID)
	if err != nil {
		return nil, err
	}

	var res []response.Role
	for _, role := range roles {
		res = append(res, response.Role{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	return res, nil
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
	// TODO: response.Userに変換
	return s.currentUsersRepository.GetCurrentUsers()
}

// 入室時間を取得
func (s *UserPortalService) GetEnteredTime(userID int) (time.Time, error) {
	return s.currentUsersRepository.GetEnteredTime(userID)
}

// model.userをresponse.userに変換
func (s *UserPortalService) cnvModelUserToResponseUser(user *model.User) *response.User {
	resUser := &response.User{
		ID:                user.ID,
		Name:              user.Name,
		Description:       user.Description,
		Barcode:           user.Barcode,
		Contact:           user.Contact,
		Remaining_entries: user.Remaining_entries,
		Registered_at:     user.Registered_at,
		Total_entries:     user.Total_entries,
		Allergy:           user.Allergy,
	}

	return resUser
}

// ユーザーが在室しているか確認
func (s *UserPortalService) IsCurrentUser(userID int) (bool, error) {
	currentUsers, err := s.currentUsersRepository.GetCurrentUsers()
	if err != nil {
		return false, err
	}

	for _, user := range currentUsers {
		if user.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

// ユーザーが在室しているか確認
func (s *UserPortalService) IsCurrentUserByBarcode(barcode string) (bool, error) {
	user, err := s.userRepository.GetUserByBarcode(barcode)
	if err != nil {
		return false, err
	}

	return s.IsCurrentUser(user.ID)
}

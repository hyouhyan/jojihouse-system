package service

import (
	"fmt"
	"jojihouse-system/api/model/request"
	"jojihouse-system/api/model/response"
	"jojihouse-system/internal/model"
	"jojihouse-system/internal/repository"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminManagementService struct {
	userRepository                *repository.UserRepository
	roleRepository                *repository.RoleRepository
	accessLogRepository           *repository.AccessLogRepository
	remainingEntriesLogRepository *repository.RemainingEntriesLogRepository
	paymentLogRepository          *repository.PaymentLogRepository
}

func NewAdminManagementService(userRepository *repository.UserRepository, roleRepository *repository.RoleRepository, accessLogRepository *repository.AccessLogRepository, remainingEntriesLogRepository *repository.RemainingEntriesLogRepository, paymentLogRepository *repository.PaymentLogRepository) *AdminManagementService {
	return &AdminManagementService{userRepository: userRepository, roleRepository: roleRepository, accessLogRepository: accessLogRepository, remainingEntriesLogRepository: remainingEntriesLogRepository, paymentLogRepository: paymentLogRepository}
}

func (s *AdminManagementService) CreateUser(req *request.CreateUser) (*response.User, error) {
	// パース的な、model.userに合わせて再構築
	user := &model.User{
		Name:              req.Name,
		Description:       req.Description,
		Barcode:           req.Barcode,
		DiscordID:         req.DiscordID,
		Remaining_entries: req.Remaining_entries,
		Allergy:           req.Allergy,
		Number:            req.Number,
	}

	// ユーザーを作成
	user, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	res := &response.User{
		ID:                user.ID,
		Name:              user.Name,
		Description:       user.Description,
		Barcode:           user.Barcode,
		DiscordID:         user.DiscordID,
		Remaining_entries: user.Remaining_entries,
		Registered_at:     user.Registered_at,
		Total_entries:     user.Total_entries,
		Allergy:           user.Allergy,
		Number:            user.Number,
	}

	log.Println("[AdminManagementService] New user created: ", *user.Name)

	return res, nil
}

func (s *AdminManagementService) UpdateUser(userID int, user *request.UpdateUser) error {
	userModel, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Name != nil {
		userModel.Name = user.Name
	}
	if user.Description != nil {
		userModel.Description = user.Description
	}
	if user.Barcode != nil {
		userModel.Barcode = user.Barcode
	}
	if user.DiscordID != nil {
		userModel.DiscordID = user.DiscordID
	}
	if user.Remaining_entries != nil {
		userModel.Remaining_entries = user.Remaining_entries
	}
	if user.Allergy != nil {
		userModel.Allergy = user.Allergy
	}
	if user.Number != nil {
		userModel.Number = user.Number
	}

	log.Println("[AdminManagementService] User updated: ", *userModel.Name)

	return s.userRepository.UpdateUser(userModel)
}

func (s *AdminManagementService) DeleteUser(userID int) error {
	// ユーザーが実在するかの確認
	_, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	return s.userRepository.DeleteUser(userID)
}

// 入場可能回数の追加
func (s *AdminManagementService) increaseRemainingEntries(userID int, count int, reason string, updatedBy string) (*primitive.ObjectID, error) {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// ユーザーが存在しない場合はエラーを返す
	if user == nil {
		return nil, model.ErrUserNotFound
	}

	// 入場可能回数 追加
	beforeCount, afterCount, err := s.userRepository.IncreaseRemainingEntries(*user.ID, count)
	if err != nil {
		return nil, err
	}

	// ログに保存
	logData := &model.RemainingEntriesLog{
		UserID:          *user.ID,
		PreviousEntries: beforeCount,
		NewEntries:      afterCount,
		Reason:          reason,
		UpdatedBy:       updatedBy,
		UpdatedAt:       time.Now(),
	}

	log.Printf("[AdminManagementService] %s's remaining entries increased %d -> %d: %s", *user.Name, beforeCount, afterCount, reason)

	return s.remainingEntriesLogRepository.CreateRemainingEntriesLog(logData)
}

func (s *AdminManagementService) decreaseRemainingEntries(userID int, count int, reason string, updatedBy string) (*primitive.ObjectID, error) {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// ユーザーが存在しない場合はエラーを返す
	if user == nil {
		return nil, model.ErrUserNotFound
	}

	// 入場可能回数を減らす
	beforeCount, afterCount, err := s.userRepository.DecreaseRemainingEntries(userID, count)
	if err != nil {
		return nil, err
	}

	// ログに保存
	logData := &model.RemainingEntriesLog{
		UserID:          *user.ID,
		PreviousEntries: beforeCount,
		NewEntries:      afterCount,
		Reason:          reason,
		UpdatedBy:       updatedBy,
		UpdatedAt:       time.Now(),
	}

	log.Printf("[AdminManagementService] %s's remaining entries decreased %d -> %d: %s", *user.Name, beforeCount, afterCount, reason)

	// ログ作成
	return s.remainingEntriesLogRepository.CreateRemainingEntriesLog(logData)
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

func (s *AdminManagementService) CreatePaymentLog(logData *model.PaymentLog) (*primitive.ObjectID, error) {
	if logData == nil {
		return nil, model.ErrInvalidPaymentLog
	}

	log.Println("[AdminManagementService] Creating payment log:", logData.Description)

	return s.paymentLogRepository.CreatePaymentLog(logData)
}

func (s *AdminManagementService) GetAllPaymentLogs(lastID string, limit int64) ([]response.PaymentLog, error) {
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

	logs, err := s.paymentLogRepository.GetAllPaymentLogs(objectID, limit)
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
		userMap[*user.ID] = *user.Name
	}

	// レスポンスデータの作成
	var responseLogs []response.PaymentLog
	for _, log := range logs {
		responseLogs = append(responseLogs, response.PaymentLog{
			ID:          log.ID.Hex(),
			UserID:      log.UserID,
			UserName:    userMap[log.UserID], // UserIDからUserNameを取得
			Time:        log.Time,
			Description: log.Description,
			Amount:      log.Amount,
			Payway:      log.Payway,
		})
	}

	return responseLogs, nil
}

func (s *AdminManagementService) GetPaymentLogByID(logID string) (*response.PaymentLog, error) {
	objectID, err := primitive.ObjectIDFromHex(logID)
	if err != nil {
		return nil, err
	}

	log, err := s.paymentLogRepository.GetPaymentLogByID(objectID)
	if err != nil {
		return nil, err
	}
	if log == nil {
		return nil, model.ErrPaymentLogNotFound
	}

	// UserIDに対応するUserNameを取得
	user, err := s.userRepository.GetUserByID(log.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, model.ErrUserNotFound
	}

	// レスポンスデータの作成
	responseLog := &response.PaymentLog{
		ID:          log.ID.Hex(),
		UserID:      log.UserID,
		UserName:    *user.Name, // UserNameをセット
		Time:        log.Time,
		Description: log.Description,
		Amount:      log.Amount,
		Payway:      log.Payway,
		IsDeleted:   log.IsDeleted,
		DeletedBy:   log.DeletedBy,
		DeletedAt:   log.DeletedAt,
	}

	return responseLog, nil
}

func (s *AdminManagementService) GetMonthlyPaymentLogs(year int, month int) (response.MonthlyPaymentLog, error) {
	monthlyLog, err := s.paymentLogRepository.GetMonthlyPaymentLogs(year, month)
	if err != nil {
		return response.MonthlyPaymentLog{}, err
	}

	// UserIDの一覧を作成
	userIDs := make([]int, len(monthlyLog.Logs))
	for i, log := range monthlyLog.Logs {
		userIDs[i] = log.UserID
	}

	// PostgreSQL から UserID に対応する UserName を取得
	users, err := s.userRepository.GetUsersByIDs(userIDs)
	if err != nil {
		return response.MonthlyPaymentLog{}, err
	}

	// UserID -> UserName のマッピング
	userMap := make(map[int]string)
	for _, user := range users {
		userMap[*user.ID] = *user.Name
	}

	// レスポンスデータの作成
	var responseLogs []response.PaymentLog
	for _, log := range monthlyLog.Logs {
		responseLogs = append(responseLogs, response.PaymentLog{
			ID:          log.ID.Hex(),
			UserID:      log.UserID,
			UserName:    userMap[log.UserID],
			Time:        log.Time,
			Description: log.Description,
			Amount:      log.Amount,
			Payway:      log.Payway,
		})
	}

	return response.MonthlyPaymentLog{
		Year:       monthlyLog.Year,
		Month:      monthlyLog.Month,
		Total:      monthlyLog.Total,
		OliveTotal: monthlyLog.OliveTotal,
		CashTotal:  monthlyLog.CashTotal,
		Logs:       responseLogs,
	}, nil
}

func (s *AdminManagementService) BuyKaisuken(userID int, receiver string, amount int, count int, payway string, description string) (*model.PaymentLog, error) {
	logDescription := fmt.Sprintf("回数券購入 %d回分 %d円", count, amount)
	if description != "" {
		logDescription = logDescription + "(" + description + ")"
	}

	remainingEntriesLogID, err := s.increaseRemainingEntries(userID, count, logDescription, receiver)
	if err != nil {
		return nil, err
	}

	// 支払いログの作成
	paymentLog := &model.PaymentLog{
		UserID:      userID,
		Amount:      amount,
		Description: logDescription,
		Payway:      payway,
	}
	paymentLogID, err := s.CreatePaymentLog(paymentLog)
	if err != nil {
		return nil, err
	}

	err = s.paymentLogRepository.LinkPaymentAndRemainingEntries(*paymentLogID, *remainingEntriesLogID)
	if err != nil {
		return nil, err
	}

	return paymentLog, nil
}

func (s *AdminManagementService) DeletePaymentLog(logID string) error {
	objectID, err := primitive.ObjectIDFromHex(logID)
	if err != nil {
		return err
	}

	// PaymentLogからRemainingEntriesLogのIDを取得
	paymentLog, err := s.paymentLogRepository.GetPaymentLogByID(objectID)
	if err != nil {
		return err
	}
	if paymentLog == nil {
		return model.ErrPaymentLogNotFound
	}

	// セーフティチェック
	if paymentLog.RemainingEntiriesLogID == nil {
		// Descriptionに"回数券購入"が含まれているか
		if strings.Contains(paymentLog.Description, "回数券購入") {
			log.Printf("[AdminManagementService] Warning: Payment log %s seems to be for ticket purchase but has no linked RemainingEntriesLogID.", logID)
			return model.ErrPaymentLogSeemsTicketPurchase
		}
	}

	// ログが14日以上前のものであれば削除不可
	if time.Since(paymentLog.Time) > 14*24*time.Hour {
		log.Print("[AdminManagementService] Warning: Payment log ", logID, " is too old to delete.", "\n", "today: ", time.Now(), " log time: ", paymentLog.Time)
		return model.ErrPaymentLogTooOldToDelete
	}

	// PaymentLogの削除
	err = s.paymentLogRepository.DeletePaymentLog(objectID)
	if err != nil {
		return err
	}
	log.Print("[AdminManagementService] Payment log deleted.\n", *paymentLog)

	// 関連RemainingEntriesログがあれば入場可能回数を戻す
	if paymentLog.RemainingEntiriesLogID != nil {
		remainintEntriesLog, err := s.remainingEntriesLogRepository.GetRemainingEntriesLogByID(*paymentLog.RemainingEntiriesLogID)
		if err != nil {
			return err
		}
		count := remainintEntriesLog.NewEntries - remainintEntriesLog.PreviousEntries

		// 入場可能回数を減らす
		_, err = s.decreaseRemainingEntries(remainintEntriesLog.UserID, count, "回数券購入の取り消しによる", "システム")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AdminManagementService) DecreaseRemainingEntries(userID int, count int, reason string, updatedBy string) (*primitive.ObjectID, error) {
	return s.decreaseRemainingEntries(userID, count, reason, updatedBy)
}

func (s *AdminManagementService) IncreaseRemainingEntries(userID int, count int, reason string, updatedBy string) (*primitive.ObjectID, error) {
	return s.increaseRemainingEntries(userID, count, reason, updatedBy)
}

package service

import (
	"fmt"
	"jojihouse-entrance-system/api/model/request"
	"jojihouse-entrance-system/api/model/response"
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
	"time"
)

type EntranceService struct {
	userRepository                *repository.UserRepository
	roleRepository                *repository.RoleRepository
	accessLogRepository           *repository.AccessLogRepository
	remainingEntriesLogRepository *repository.RemainingEntriesLogRepository
	currentUsersRepository        *repository.CurrentUsersRepository
	discordNoticeRepository       *repository.DiscordNoticeRepository
}

func NewEntranceService(
	userRepository *repository.UserRepository,
	roleRepository *repository.RoleRepository,
	accessLogRepository *repository.AccessLogRepository,
	remainingEntriesLogRepository *repository.RemainingEntriesLogRepository,
	currentUsersRepository *repository.CurrentUsersRepository,
	discordNoticeRepository *repository.DiscordNoticeRepository,
) *EntranceService {
	return &EntranceService{
		userRepository:                userRepository,
		roleRepository:                roleRepository,
		accessLogRepository:           accessLogRepository,
		remainingEntriesLogRepository: remainingEntriesLogRepository,
		currentUsersRepository:        currentUsersRepository,
		discordNoticeRepository:       discordNoticeRepository,
	}
}

// 入場したときの処理
func (s *EntranceService) EnterUser(barcode string) (response.Entrance, error) {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.userRepository.GetUserByBarcode(barcode)
	if err != nil {
		return response.Entrance{}, err
	}

	// 入場ログ作成
	err = s.accessLogRepository.CreateEntryAccessLog(*user.ID)
	if err != nil {
		return response.Entrance{}, err
	}

	// 在室ユーザーに追加
	err = s.currentUsersRepository.AddUserToCurrentUsers(*user.ID)
	if err != nil {
		return response.Entrance{}, err
	}

	isDecreaseTarget := true
	// ハウス管理者か
	isHouseAdmin, err := s.roleRepository.IsHouseAdmin(*user.ID)
	if err != nil {
		return response.Entrance{}, fmt.Errorf("failed to check if the user is a house admin: %v", err)
	}
	if isHouseAdmin {
		isDecreaseTarget = false
	}

	// 最後に「入場可能回数を消費した」入場を取得
	lastRemainingLog, err := s.remainingEntriesLogRepository.GetLastRemainingEntriesLogByUserID(*user.ID)
	if err != nil {
		return response.Entrance{}, err
	}

	// DEBUG
	fmt.Printf("Last Entries Date: %s, Current Date: %s\n", lastRemainingLog.UpdatedAt.Format("2006-01-02"), time.Now().Format("2006-01-02"))

	// ログの日が今日なら同日再入場
	if isSameDate(lastRemainingLog.UpdatedAt, time.Now()) {
		isDecreaseTarget = false
	}

	if isDecreaseTarget {
		// 残り回数を減らす
		beforeCount, afterCount, err := s.userRepository.DecreaseRemainingEntries(*user.ID, 1)
		if err != nil {
			return response.Entrance{}, err
		}
		// ログ保存

		log := &model.RemainingEntriesLog{
			UserID:          *user.ID,
			PreviousEntries: beforeCount,
			NewEntries:      afterCount,
			Reason:          "ハウス入場のため",
			UpdatedBy:       "システム",
		}

		// ログ作成
		err = s.remainingEntriesLogRepository.CreateRemainingEntriesLog(log)
		if err != nil {
			return response.Entrance{}, err
		}

		// Go側にも反映
		*user.Remaining_entries = afterCount
	}

	// 入場回数を増やす
	err = s.userRepository.IncreaseTotalEntries(*user.ID)
	if err != nil {
		return response.Entrance{}, err
	}
	// Go側にも反映
	*user.Total_entries = *user.Total_entries + 1

	// Discordに通知
	go s.discordNoticeRepository.NoticeEntry(*user.Name)

	// Response作成
	response := response.Entrance{
		UserID:            *user.ID,
		UserName:          *user.Name,
		Time:              time.Now(),
		AccessType:        "entry",
		Remaining_entries: *user.Remaining_entries,
		Number:            user.Number,
		Total_entries:     *user.Total_entries,
	}

	return response, nil
}

// 退場したときの処理
func (s *EntranceService) ExitUser(barcode string) (response.Entrance, error) {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.userRepository.GetUserByBarcode(barcode)
	if err != nil {
		return response.Entrance{}, err
	}

	// 退場ログ作成
	err = s.accessLogRepository.CreateExitAccessLog(*user.ID)
	if err != nil {
		return response.Entrance{}, err
	}

	// 在室ユーザーから削除
	err = s.currentUsersRepository.DeleteUserToCurrentUsers(*user.ID)
	if err != nil {
		return response.Entrance{}, err
	}

	// 日をまたいでいないか確認
	// 最後に「入場可能回数を消費した」入場を取得
	lastRemainingLog, err := s.remainingEntriesLogRepository.GetLastRemainingEntriesLogByUserID(*user.ID)
	if err != nil {
		return response.Entrance{}, err
	}

	// DEBUG
	fmt.Printf("Last Entries Date: %s, Current Date: %s\n", lastRemainingLog.UpdatedAt.Format("2006-01-02"), time.Now().Format("2006-01-02"))

	// ログ日から今日までの時間差を確認
	if !isSameDate(lastRemainingLog.UpdatedAt, time.Now()) {
		// 何日経過したかの計算
		daysPassed := int(time.Since(lastRemainingLog.UpdatedAt).Hours() / 24)
		if daysPassed == 0 {
			fmt.Println("起こり得ないエラー: 日を跨いでいるのに経過日数が0")
		}

		// 残り回数を減らす
		beforeCount, afterCount, err := s.userRepository.DecreaseRemainingEntries(*user.ID, 1)
		if err != nil {
			return response.Entrance{}, err
		}
		// ログ保存

		log := &model.RemainingEntriesLog{
			UserID:          *user.ID,
			PreviousEntries: beforeCount,
			NewEntries:      afterCount,
			Reason:          "ハウス入場(日跨ぎ)のため",
			UpdatedBy:       "システム",
		}

		// ログ作成
		err = s.remainingEntriesLogRepository.CreateRemainingEntriesLog(log)
		if err != nil {
			return response.Entrance{}, err
		}

		// Go側にも反映
		*user.Remaining_entries = afterCount
	}

	// Discordに通知
	go s.discordNoticeRepository.NoticeExit(*user.Name)

	// Response作成
	response := response.Entrance{
		UserID:            *user.ID,
		UserName:          *user.Name,
		Time:              time.Now(),
		AccessType:        "exit",
		Remaining_entries: *user.Remaining_entries,
		Number:            user.Number,
		Total_entries:     *user.Total_entries,
	}

	return response, nil
}

func isSameDate(a, b time.Time) bool {
	aDate := a.Truncate(24 * time.Hour)
	bDate := b.Truncate(24 * time.Hour)
	return aDate.Equal(bDate)
}

// 入退室の修正
func (s *EntranceService) UpdateAccessLog(log *model.AccessLog) error {
	// 既存のログを更新
	err := s.accessLogRepository.UpdateAccessLog(log)
	if err != nil {
		return err
	}

	// 最終アクセスログを取得
	lastLog, err := s.accessLogRepository.GetLastAccessLogByUserID(log.UserID)
	if err != nil {
		return err
	}

	// 在室ユーザーを取得
	currentUsers, err := s.currentUsersRepository.GetCurrentUsers()
	if err != nil {
		return err
	}
	// ユーザーが在室中か確認
	var isCurrentUser bool
	for _, user := range currentUsers {
		if user.UserID == log.UserID {
			isCurrentUser = true
			break
		}
	}

	// 最終アクセスを元に在室ユーザーを更新
	if lastLog.AccessType == "exit" && isCurrentUser {
		// 退場ログが最新で、在室ユーザーにいる場合は、在室ユーザーから削除
		err = s.currentUsersRepository.DeleteUserToCurrentUsers(log.UserID)
		if err != nil {
			return err
		}
	}
	if lastLog.AccessType == "entry" && !isCurrentUser {
		// 入場ログが最新で、在室ユーザーにいない場合は、在室ユーザーに追加
		err = s.currentUsersRepository.AddUserToCurrentUsers(log.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *EntranceService) CreateFixedAccessLog(req *request.FixedAccessLog) error {
	// ユーザー情報を取得(存在するかの確認)
	var user *model.User
	var err error

	if req.UserID != nil {
		user, err = s.userRepository.GetUserByID(*req.UserID)
		if err != nil {
			return err
		}
	} else if req.Barcode != nil {
		user, err = s.userRepository.GetUserByBarcode(*req.Barcode)
		if err != nil {
			return err
		}
	} else if req.Number != nil {
		user, err = s.userRepository.GetUserByNumber(*req.Number)
		if err != nil {
			return err
		}
	}

	log := &model.AccessLog{
		UserID:     *user.ID,
		Time:       *req.Time,
		AccessType: *req.AccessType,
	}

	// ログを作成
	err = s.accessLogRepository.CreateAccessLog(log)
	if err != nil {
		return err
	}

	// 最終アクセスログを取得
	lastLog, err := s.accessLogRepository.GetLastAccessLogByUserID(log.UserID)
	if err != nil {
		return err
	}

	// 在室ユーザーを取得
	currentUsers, err := s.currentUsersRepository.GetCurrentUsers()
	if err != nil {
		return err
	}
	// ユーザーが在室中か確認
	var isCurrentUser bool
	for _, user := range currentUsers {
		if user.UserID == log.UserID {
			isCurrentUser = true
			break
		}
	}

	// 最終アクセスを元に在室ユーザーを更新
	if lastLog.AccessType == "exit" && isCurrentUser {
		// 退場ログが最新で、在室ユーザーにいる場合は、在室ユーザーから削除
		err = s.currentUsersRepository.DeleteUserToCurrentUsers(log.UserID)
		if err != nil {
			return err
		}
	}
	if lastLog.AccessType == "entry" && !isCurrentUser {
		// 入場ログが最新で、在室ユーザーにいない場合は、在室ユーザーに追加
		err = s.currentUsersRepository.AddUserToCurrentUsers(log.UserID)
		if err != nil {
			return err
		}
	}

	// 入場時は総入場回数を増やす
	if lastLog.AccessType == "entry" {
		err = s.userRepository.IncreaseTotalEntries(log.UserID)
		if err != nil {
			return err
		}

		// その日のうちにremaining_entries_logの減少があるか確認
		lastRemainingLog, err := s.remainingEntriesLogRepository.GetLastRemainingEntriesLogByUserID(log.UserID)
		if err != nil {
			return err
		}

		// ハウス管理者か
		isHouseAdmin, err := s.roleRepository.IsHouseAdmin(log.UserID)
		if err != nil {
			return err
		}

		// 指定日とログの日が同じなら、減少しない
		if isSameDate(lastRemainingLog.UpdatedAt, *req.Time) || isHouseAdmin {
			return nil
		} else {
			// 残り回数を減らす
			beforeCount, afterCount, err := s.userRepository.DecreaseRemainingEntries(log.UserID, 1)
			if err != nil {
				return err
			}
			// ログ保存

			log := &model.RemainingEntriesLog{
				UserID:          log.UserID,
				PreviousEntries: beforeCount,
				NewEntries:      afterCount,
				Reason:          "ハウス入場(修正)のため",
				UpdatedBy:       "システム",
				UpdatedAt:       *req.Time,
			}

			// ログ作成
			err = s.remainingEntriesLogRepository.CreateFixedRemainingEntriesLog(log)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

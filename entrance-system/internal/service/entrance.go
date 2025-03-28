package service

import (
	"jojihouse-entrance-system/internal/model"
	"log"
	"time"
)

type EntranceService struct {
	userService                *UserService
	roleService                *RoleService
	accessLogService           *AccessLogService
	remainingEntriesLogService *RemainingEntriesLogService
}

func NewEntranceService(userService *UserService, roleService *RoleService, accessLogService *AccessLogService, remainingEntriesLogService *RemainingEntriesLogService) *EntranceService {
	return &EntranceService{userService: userService, roleService: roleService, accessLogService: accessLogService, remainingEntriesLogService: remainingEntriesLogService}
}

// 入場したときの処理
func (s *EntranceService) EnterUser(barcode string) error {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.userService.GetUserByBarcode(barcode)
	if err != nil {
		return err
	}

	// 入場ログ作成
	err = s.accessLogService.CreateEntryAccessLog(user.ID)
	if err != nil {
		return err
	}

	isDecreaseTarget := true
	// ハウス管理者か
	isHouseAdmin, err := s.roleService.IsHouseAdmin(user.ID)
	if err != nil {
		log.Fatalf("Failed to check if the user is a house admin: %v", err)
	}
	if isHouseAdmin {
		isDecreaseTarget = false
	}

	// 最後に「入場可能回数を消費した」入場を取得
	lastRemainingLog, err := s.remainingEntriesLogService.GetLastRemainingEntriesLogByUserID(user.ID)
	if err != nil {
		return err
		// mongo: no documents in resultが返される
		// 登録されたばかりのユーザーや、ハウス管理者の場合は正常な動作
		// エラーをreturnしないか、このエラーだけ特例で許すか
	}

	// ログの日が今日なら同日再入場
	if isSameDate(lastRemainingLog.UpdatedAt, time.Now()) {
		isDecreaseTarget = false
	}

	// 取得したログと同日の入場かどうか

	if isDecreaseTarget {
		// 残り回数を減らす
		err = s.userService.DecreaseRemainingEntries(user.ID)
		if err != nil {
			return err
		}
		// ログ保存
		// 変更前残り回数
		prevRemain := user.Remaining_entries
		user, err := s.userService.GetUserByID(user.ID)
		if err != nil {
			return err
		}
		// 変更後残り回数
		newRemain := user.Remaining_entries

		log := &model.RemainingEntriesLog{
			UserID:          user.ID,
			PreviousEntries: prevRemain,
			NewEntries:      newRemain,
			Reason:          "ハウス入場のため",
			UpdatedBy:       "システム",
		}

		// ログ作成
		s.remainingEntriesLogService.CreateRemainingEntriesLog(log)
	}

	return nil
}

// 退場したときの処理
func (s *EntranceService) ExitUser(barcode string) error {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.userService.GetUserByBarcode(barcode)
	if err != nil {
		return err
	}

	// 退場ログ作成
	err = s.accessLogService.CreateExitAccessLog(user.ID)
	if err != nil {
		return err
	}

	// 入場回数を増やす
	err = s.userService.IncreaseTotalEntries(user.ID)
	if err != nil {
		return err
	}

	return nil
}

func isSameDate(a, b time.Time) bool {
	aDate := a.Truncate(24 * time.Hour)
	bDate := b.Truncate(24 * time.Hour)
	return aDate.Equal(bDate)
}

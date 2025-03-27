package service

import (
	"jojihouse-entrance-system/internal/model"
	"log"
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

	// TODO: 同日の再入場か

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

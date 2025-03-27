package service

import (
	"log"
)

type EntranceService struct {
	userService      *UserService
	roleService      *RoleService
	accessLogService *AccessLogService
}

func NewEntranceService(userService *UserService, roleService *RoleService, accessLogService *AccessLogService) *EntranceService {
	return &EntranceService{userService: userService, roleService: roleService, accessLogService: accessLogService}
}

// 入場したときの処理
func (s *EntranceService) EnterUser(barcode string) error {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.userService.GetUserByBarcode(barcode)
	if err != nil {
		return err
	}

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

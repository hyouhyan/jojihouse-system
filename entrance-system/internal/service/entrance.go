package service

type EntranceService struct {
	userService *UserService
	roleService *RoleService
}

func NewEntranceService(userService *UserService, roleService *RoleService) *EntranceService {
	return &EntranceService{userService: userService, roleService: roleService}
}

// 入場したときの処理
func (s *EntranceService) EnterUser(barcode string) error {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.userService.GetUserByBarcode(barcode)
	if err != nil {
		return err
	}

	// TODO: ログの生成

	// TODO: ユーザーが入場可能回数を減らす対象かの確認
	// ハウス管理者とか、同日の再入場とか

	// 残り回数を減らす
	err = s.userService.DecreaseRemainingEntries(user.ID)
	if err != nil {
		return err
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

	// TODO: ログの生成

	// 入場回数を増やす
	err = s.userService.IncreaseTotalEntries(user.ID)
	if err != nil {
		return err
	}

	return nil
}

package service

import (
	"fmt"
	"jojihouse-system/api/model/response"
	"jojihouse-system/internal/model"
	"jojihouse-system/internal/repository"
	"log"
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
	log.Println("[EntranceService] Enter Requested: ", barcode)

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
		log.Printf("[EntranceService] %s is a house admin", *user.Name)
	}

	// 最後に「入場可能回数を消費した」入場を取得
	lastRemainingLog, err := s.remainingEntriesLogRepository.GetLastDecreaseRemainingEntriesLogByUserID(*user.ID)
	if err != nil {
		return response.Entrance{}, err
	}

	// 計算のためにTimezoneをlocalに変換
	lastDate := lastRemainingLog.UpdatedAt.In(time.Local)
	currentDate := time.Now()

	// ログの日が今日なら同日再入場
	if s.isSameDate(lastDate, currentDate) {
		isDecreaseTarget = false
		log.Printf("[EntranceService] %s re-entered on the same day", *user.Name)
	}

	if isDecreaseTarget {
		// 残り回数を減らす
		beforeCount, afterCount, err := s.userRepository.DecreaseRemainingEntries(*user.ID, 1)
		if err != nil {
			return response.Entrance{}, err
		}

		// ログ保存
		logData := &model.RemainingEntriesLog{
			UserID:          *user.ID,
			PreviousEntries: beforeCount,
			NewEntries:      afterCount,
			Reason:          "ハウス入場のため",
			UpdatedBy:       "システム",
		}

		// ログ作成
		_, err = s.remainingEntriesLogRepository.CreateRemainingEntriesLog(logData)
		if err != nil {
			return response.Entrance{}, err
		}

		// Go側にも反映
		*user.Remaining_entries = afterCount

		log.Printf("[EntranceService] %s's remaining entries decreased %d -> %d", *user.Name, beforeCount, afterCount)
	}

	// 入場回数を増やす
	err = s.userRepository.IncreaseTotalEntries(*user.ID)
	if err != nil {
		return response.Entrance{}, err
	}
	// Go側にも反映
	*user.Total_entries = *user.Total_entries + 1

	// Logに出力
	log.Printf("[EntranceService] %s entered. Barcode: %s, Remaining entries: %d, Total entries: %d", *user.Name, *user.Barcode, *user.Remaining_entries, *user.Total_entries)

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
	log.Println("[EntranceService] Exit Requested: ", barcode)

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
	lastRemainingLog, err := s.remainingEntriesLogRepository.GetLastDecreaseRemainingEntriesLogByUserID(*user.ID)
	if err != nil {
		return response.Entrance{}, err
	}

	isHouseAdmin, err := s.roleRepository.IsHouseAdmin(*user.ID)
	if err != nil {
		return response.Entrance{}, fmt.Errorf("failed to check if the user is a house admin: %v", err)
	}

	// 計算のためにTimezoneをlocalに変換
	lastDate := lastRemainingLog.UpdatedAt.In(time.Local)
	currentDate := time.Now()

	// ログ日から今日までの時間差を確認
	if !s.isSameDate(lastDate, currentDate) && !isHouseAdmin {
		// 何日経過したかの計算
		daysPassed := s.getPassedDays(lastDate, currentDate)
		if daysPassed == 0 {
			log.Println("[EntranceService] 起こり得ないエラー: 日を跨いでいるのに経過日数が0")

			log.Println("[EntranceService] lastDate: ", lastDate)
			log.Println("[EntranceService] currentDate: ", currentDate)
		}

		log.Printf("[EntranceService] Days Passed: %d\n", daysPassed)

		// 残り回数を減らす
		beforeCount, afterCount, err := s.userRepository.DecreaseRemainingEntries(*user.ID, daysPassed)
		if err != nil {
			return response.Entrance{}, err
		}

		// ログ保存
		logData := &model.RemainingEntriesLog{
			UserID:          *user.ID,
			PreviousEntries: beforeCount,
			NewEntries:      afterCount,
			Reason:          "日付を跨いだハウス利用のため",
			UpdatedBy:       "システム",
		}

		// ログ作成
		_, err = s.remainingEntriesLogRepository.CreateRemainingEntriesLog(logData)
		if err != nil {
			return response.Entrance{}, err
		}

		// Go側にも反映
		*user.Remaining_entries = afterCount

		log.Printf("[EntranceService] %s's remaining entries decreased %d -> %d due to date crossing", *user.Name, beforeCount, afterCount)
	}

	// Logに出力
	log.Printf("[EntranceService] %s exited. Barcode: %s, Remaining entries: %d, Total entries: %d", *user.Name, *user.Barcode, *user.Remaining_entries, *user.Total_entries)

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

func (s *EntranceService) isSameDate(a, b time.Time) bool {
	// bのtimezoneをaのtimezoneに合わせる
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}

	aDate := s.cnvTo00Time(a)
	bDate := s.cnvTo00Time(b)

	return aDate.Equal(bDate)
}

func (s *EntranceService) cnvTo00Time(t time.Time) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0, t.Location())
}

func (s *EntranceService) getPassedDays(targetDate, currentDate time.Time) int {
	// aとbのtimezoneを揃える
	if targetDate.Location() != currentDate.Location() {
		currentDate = currentDate.In(targetDate.Location())
	}

	// 00:00:00どうしで比較
	targetDate = s.cnvTo00Time(targetDate)
	currentDate = s.cnvTo00Time(currentDate)

	// 日数の差を計算
	daysPassed := int(currentDate.Sub(targetDate).Hours() / 24)

	return daysPassed
}

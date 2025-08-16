package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jojihouse-management-system/internal/config"
	"log"
	"net/http"
	"os"
	"time"
)

// Discord Webhookに送信するJSONデータ構造を定義します
type WebhookPayload struct {
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Embeds    []Embed `json:"embeds"`
}

type Embed struct {
	Title  string `json:"title"`
	Footer Footer `json:"footer"`
	Color  int    `json:"color"`
}

type Footer struct {
	Text string `json:"text"`
}

// 日本語の曜日スライス
var weekdays = []string{"日", "月", "火", "水", "木", "金", "土"}

type DiscordNoticeRepository struct{}

func NewDiscordNoticeRepository() *DiscordNoticeRepository {
	return &DiscordNoticeRepository{}
}

func (r *DiscordNoticeRepository) NoticeEntry(userName string) {
	err := r.noticeAccess(userName, "入室")
	if err != nil {
		log.Printf("Error sending entry notice: %v\n", err)
	} else {
		log.Println("Entry notice sent successfully")
	}
}

func (r *DiscordNoticeRepository) NoticeExit(userName string) {
	err := r.noticeAccess(userName, "退室")
	if err != nil {
		log.Printf("Error sending exit notice: %v\n", err)
	} else {
		log.Println("Exit notice sent successfully")
	}
}

func (r *DiscordNoticeRepository) noticeAccess(userName string, accessType string) error {
	config.Env_load()

	WEBHOOK_URL := os.Getenv("WEBHOOK_URL")
	WEBHOOK_USERNAME := os.Getenv("WEBHOOK_USERNAME")
	WEBHOOK_AVATAR_URL := os.Getenv("WEBHOOK_AVATAR_URL")

	colorGreen := 0x2ECC71 // 10進数: 3066993
	colorRed := 0xE74C3C   // 10進数: 15158332

	var color int
	if accessType == "入室" {
		color = colorGreen
	} else {
		color = colorRed
	}

	// 現在時刻を取得し、フッター用のテキストを生成
	now := time.Now()
	footerText := fmt.Sprintf(
		"%s(%s) %s",
		now.Format("2006年01月02日"), // YYYY年MM月DD日
		weekdays[now.Weekday()],   // (曜日)
		now.Format("15:04:05"),    // HH:MM:SS
	)

	// 送信するデータ（ペイロード）を作成
	payload := WebhookPayload{
		Username:  WEBHOOK_USERNAME,
		AvatarURL: WEBHOOK_AVATAR_URL,
		Embeds: []Embed{
			{
				Title: fmt.Sprintf("%sが%sしました", userName, accessType),
				Footer: Footer{
					Text: footerText,
				},
				Color: color,
			},
		},
	}

	// ペイロードをJSON形式のバイト配列に変換（マーシャリング）
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("JSONのマーシャリングに失敗しました: %v", err)
	}

	// HTTP POSTリクエストを作成して送信
	req, err := http.NewRequest("POST", WEBHOOK_URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗しました: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("リクエストの送信に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスを表示
	log.Println("Response Status:", resp.Status)

	return nil
}

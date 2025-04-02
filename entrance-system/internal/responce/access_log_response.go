package response

import (
	"time"
)

// フロントエンドに返すためのログデータ
type AccessLogResponse struct {
	ID         string    `json:"id"`
	UserID     int       `json:"user_id"`
	UserName   string    `json:"user_name"`
	Time       time.Time `json:"time"`
	AccessType string    `json:"access_type"`
}

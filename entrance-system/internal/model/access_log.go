package model

import (
	"time"
)

type AccessLog struct {
	UserID     int       `json:"user_id"`
	Time       time.Time `json:"time"`
	AccessType string    `json:"access_type"`
}

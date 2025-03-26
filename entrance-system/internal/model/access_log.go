package model

import (
	"time"
)

type AccessLog struct {
	ID         int       `json:"id" gorm:"primary_key"`
	UserID     int       `json:"user_id"`
	Time       time.Time `json:"time"`
	AccessType string    `json:"access_type"`
}

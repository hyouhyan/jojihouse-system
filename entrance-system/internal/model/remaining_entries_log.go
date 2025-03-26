package model

import (
	"time"
)

type RemainingEntriesLog struct {
	ID              int       `json:"id" gorm:"primary_key"`
	UserID          int       `json:"user_id"`
	PreviousEntries int       `json:"previous_entries"`
	NewEntries      int       `json:"new_entries"`
	Reason          string    `json:"reason"`
	UpdatedBy       string    `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}

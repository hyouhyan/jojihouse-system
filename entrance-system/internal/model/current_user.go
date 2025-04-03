package model

import "time"

type CurrentUser struct {
	UserID    int       `json:"user_id"`
	Name      string    `json:"user_name"`
	EnteredAt time.Time `json:"entered_at"`
}

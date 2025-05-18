package model

import "time"

type CurrentUser struct {
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"user_name" db:"name"`
	EnteredAt time.Time `json:"entered_at" db:"entered_at"`
	Allergy   string    `json:"allergy" db:"allergy"`
}

package response

import (
	"time"
)

type Entrance struct {
	UserID     int       `json:"user_id"`
	UserName   string    `json:"user_name"`
	Time       time.Time `json:"time"`
	AccessType string    `json:"access_type"`
}

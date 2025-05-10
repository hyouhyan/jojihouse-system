package request

import "time"

type Entrance struct {
	Barcode string `json:"barcode" binding:"required"`
	Type    string `json:"type" binding:"required,oneof=entry exit auto"`
}

type FixedAccessLog struct {
	Barcode    *string    `json:"barcode"`
	UserID     *int       `json:"user_id"`
	Number     *int       `json:"number"`
	Time       *time.Time `json:"time" binding:"required"`
	AccessType *string    `json:"access_type" binding:"required,oneof=entry exit"`
}

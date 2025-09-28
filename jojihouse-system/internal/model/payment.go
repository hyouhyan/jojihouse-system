package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidPaymentLog = errors.New("invalid payment log")

type PaymentLog struct {
	ID                     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID                 int                `json:"user_id" bson:"user_id"`
	Time                   time.Time          `json:"time" bson:"time"`
	Description            string             `json:"description" bson:"description"`
	Amount                 int                `json:"amount" bson:"amount"`
	Payway                 string             `json:"payway" bson:"payway"`
	RemainingEntiriesLogID int                `json:"remaining_entries_log_id" bson:"remaining_entries_log_id"`
}

type MonthlyPaymentLog struct {
	Year       int          `json:"year" bson:"year"`
	Month      int          `json:"month" bson:"month"`
	Total      int          `json:"total" bson:"total"`
	OliveTotal int          `json:"olive_total" bson:"olive_total"`
	CashTotal  int          `json:"cash_total" bson:"cash_total"`
	Logs       []PaymentLog `json:"logs" bson:"logs"`
}

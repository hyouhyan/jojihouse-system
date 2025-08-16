package response

import (
	"time"
)

type PaymentLog struct {
	ID          string    `json:"id"`
	UserID      int       `json:"user_id"`
	UserName    string    `json:"user_name"`
	Time        time.Time `json:"time"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Payway      string    `json:"payway"`
}

type MonthlyPaymentLog struct {
	Year       int          `json:"year"`
	Month      int          `json:"month"`
	Total      int          `json:"total"`
	OliveTotal int          `json:"olive_total"`
	CashTotal  int          `json:"cash_total"`
	Logs       []PaymentLog `json:"logs"`
}

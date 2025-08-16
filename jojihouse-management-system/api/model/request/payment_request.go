package request

type Payment struct {
	UserID      int    `json:"user_id" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	Description string `json:"description" binding:"required"`
	Payway      string `json:"payway" binding:"required,oneof=olive cash"`
}

package request

type BuyKaisuken struct {
	UserID   int    `json:"user_id" binding:"required"`
	Receiver string `json:"receiver" binding:"required"`
	Amount   int    `json:"amount" binding:"required"`
	Count    int    `json:"count" binding:"required"`
}

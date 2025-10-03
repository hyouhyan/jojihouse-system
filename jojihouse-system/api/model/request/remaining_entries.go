package request

type ChangeRemainingEntries struct {
	Delta     int    `json:"delta" binding:"required"`
	Readon    string `json:"reason" binding:"required"`
	UpdatedBy string `json:"updated_by" binding:"required"`
}

package request

type CreateUserRequest struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	Barcode           string `json:"barcode" binding:"required"`
	Contact           string `json:"contact"`
	Remaining_entries int    `json:"remaining_entries"`
}

package request

type CreateUser struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	Barcode           string `json:"barcode" binding:"required"`
	Contact           string `json:"contact"`
	Remaining_entries int    `json:"remaining_entries"`
}

type UpdateUser struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Barcode           string `json:"barcode"`
	Contact           string `json:"contact"`
	Remaining_entries int    `json:"remaining_entries"`
}

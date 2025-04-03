package request

type CreateUser struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	Barcode           string `json:"barcode" binding:"required"`
	Contact           string `json:"contact"`
	Remaining_entries int    `json:"remaining_entries"`
}

type UpdateUser struct {
	Name              *string `json:"name,omitempty"`
	Description       *string `json:"description,omitempty"`
	Barcode           *string `json:"barcode,omitempty"`
	Contact           *string `json:"contact,omitempty"`
	Remaining_entries *int    `json:"remaining_entries,omitempty"`
}

package request

type CreateUser struct {
	Name              *string `json:"name" binding:"required"`
	Description       *string `json:"description"`
	Barcode           *string `json:"barcode"`
	Contact           *string `json:"contact"`
	Remaining_entries *int    `json:"remaining_entries"`
	Allergy           *string `json:"allergy"`
	Number            *int    `json:"number"`
}

type UpdateUser struct {
	Name              *string `json:"name,omitempty"`
	Description       *string `json:"description,omitempty"`
	Barcode           *string `json:"barcode,omitempty"`
	Contact           *string `json:"contact,omitempty"`
	Remaining_entries *int    `json:"remaining_entries,omitempty"`
	Allergy           *string `json:"allergy,omitempty"`
	Number            *int    `json:"number,omitempty"`
}

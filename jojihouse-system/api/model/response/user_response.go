package response

import (
	"time"
)

type User struct {
	ID                *int       `json:"id"`
	Name              *string    `json:"name"`
	Description       *string    `json:"description"`
	Barcode           *string    `json:"barcode"`
	DiscordID         *int       `json:"discord_id"`
	Remaining_entries *int       `json:"remaining_entries"`
	Registered_at     *time.Time `json:"registered_at"`
	Total_entries     *int       `json:"total_entries"`
	Allergy           *string    `json:"allergy"`
	Number            *int       `json:"number"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

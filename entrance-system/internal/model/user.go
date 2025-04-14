package model

import "time"

type User struct {
	ID                int       `db:"id"`
	Name              string    `db:"name"`
	Description       string    `db:"description"`
	Barcode           string    `db:"barcode"`
	Contact           string    `db:"contact"`
	Remaining_entries int       `db:"remaining_entries"`
	Registered_at     time.Time `db:"registered_at"`
	Total_entries     int       `db:"total_entries"`
	Allergy           string    `db:"allergy"`
}

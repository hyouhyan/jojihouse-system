package response

import "time"

type Logs struct {
	RemainingEntriesLog []RemainingEntriesLog `json:"remaining_entries_log"`
}

type RemainingEntriesLog struct {
	ID              string    `json:"id"`
	UserID          int       `json:"user_id"`
	PreviousEntries int       `json:"previous_entries"`
	NewEntries      int       `json:"new_entries"`
	Reason          string    `json:"reason"`
	UpdatedBy       string    `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RemainingEntriesLog struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          int                `json:"user_id"`
	PreviousEntries int                `json:"previous_entries"`
	NewEntries      int                `json:"new_entries"`
	Reason          string             `json:"reason"`
	UpdatedBy       string             `json:"updated_by"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

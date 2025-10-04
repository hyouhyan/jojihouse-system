package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrRemainingEntriesLogNotFound = errors.New("remaining entries log not found")

type RemainingEntriesLog struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          int                `json:"user_id" bson:"user_id"`
	PreviousEntries int                `json:"previous_entries" bson:"previous_entries"`
	NewEntries      int                `json:"new_entries" bson:"new_entries"`
	Reason          string             `json:"reason" bson:"reason"`
	UpdatedBy       string             `json:"updated_by" bson:"updated_by"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessLog struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     int                `json:"user_id" bson:"user_id"`
	Time       time.Time          `json:"time" bson:"time"`
	AccessType string             `json:"access_type" bson:"access_type"`
}

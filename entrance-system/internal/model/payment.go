package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentLog struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      int                `json:"user_id" bson:"user_id"`
	Time        time.Time          `json:"time" bson:"time"`
	Description string             `json:"description" bson:"description"`
	Amount      float64            `json:"amount" bson:"amount"`
	Payway      string             `json:"payway" bson:"payway"`
}

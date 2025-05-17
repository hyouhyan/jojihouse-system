package repository

import (
	"context"
	"jojihouse-entrance-system/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaymentLogRepository struct {
	db *mongo.Database
}

func NewPaymentLogRepository(db *mongo.Database) *PaymentLogRepository {
	return &PaymentLogRepository{db: db}
}

func (r *PaymentLogRepository) CreatePaymentLog(log *model.PaymentLog) error {
	log.ID = primitive.NilObjectID
	log.Time = time.Now()

	_, err := r.db.Collection("payment_log").InsertOne(context.Background(), log)
	if err != nil {
		return err
	}
	return nil
}

func (r *PaymentLogRepository) GetAllPaymentLogs(lastID primitive.ObjectID, limit int64) ([]model.PaymentLog, error) {
	var logs []model.PaymentLog
	filter := bson.D{}
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "_id", Value: -1}})
	opts.SetLimit(limit)
	if !lastID.IsZero() {
		filter = bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "$lt", Value: lastID},
			}},
		}
	}
	cursor, err := r.db.Collection("payment_log").Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var log model.PaymentLog
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return logs, nil
}

package repository

import (
	"context"

	"jojihouse-entrance-system/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccessLogRepository struct {
	db *mongo.Database
}

func NewLogRepository(db *mongo.Database) *AccessLogRepository {
	return &AccessLogRepository{db: db}
}

func (r *AccessLogRepository) CreateAccessLog(log *model.AccessLog) error {
	_, err := r.db.Collection("access_log").InsertOne(context.Background(), log)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccessLogRepository) GetAccessLogs() ([]model.AccessLog, error) {
	var logs []model.AccessLog
	cursor, err := r.db.Collection("access_log").Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var log model.AccessLog
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *AccessLogRepository) GetAccessLogsByUserID(userID int) ([]model.AccessLog, error) {
	var logs []model.AccessLog
	cursor, err := r.db.Collection("access_log").Find(context.Background(), map[string]int{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var log model.AccessLog
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

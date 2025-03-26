package repository

import (
	"context"

	"jojihouse-entrance-system/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type LogRepository struct {
	db *mongo.Database
}

func NewLogRepository(db *mongo.Database) *LogRepository {
	return &LogRepository{db: db}
}

func (r *LogRepository) CreateAccessLog(log *model.AccessLog) error {
	_, err := r.db.Collection("access_log").InsertOne(context.Background(), log)
	if err != nil {
		return err
	}
	return nil
}

func (r *LogRepository) GetAccessLogs() ([]model.AccessLog, error) {
	var logs []model.AccessLog
	cursor, err := r.db.Collection("access_log").Find(context.Background(), nil)
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

func (r *LogRepository) GetAccessLogsByUserID(userID int) ([]model.AccessLog, error) {
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

func (r *LogRepository) CreateRemainingEntriesLog(log *model.RemainingEntriesLog) error {
	_, err := r.db.Collection("remaining_entries_log").InsertOne(context.Background(), log)
	if err != nil {
		return err
	}
	return nil
}

func (r *LogRepository) GetRemainingEntriesLogs() ([]model.RemainingEntriesLog, error) {
	var logs []model.RemainingEntriesLog
	cursor, err := r.db.Collection("remaining_entries_log").Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var log model.RemainingEntriesLog
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *LogRepository) GetRemainingEntriesLogsByUserID(userID int) ([]model.RemainingEntriesLog, error) {
	var logs []model.RemainingEntriesLog
	cursor, err := r.db.Collection("remaining_entries_log").Find(context.Background(), map[string]int{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var log model.RemainingEntriesLog
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

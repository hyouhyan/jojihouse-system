package repository

import (
	"context"

	"jojihouse-entrance-system/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type RemainingEntriesLogRepository struct {
	db *mongo.Database
}

func NewRemainingEntriesLogRepository(db *mongo.Database) *RemainingEntriesLogRepository {
	return &RemainingEntriesLogRepository{db: db}
}

func (r *RemainingEntriesLogRepository) CreateRemainingEntriesLog(log *model.RemainingEntriesLog) error {
	_, err := r.db.Collection("remaining_entries_log").InsertOne(context.Background(), log)
	if err != nil {
		return err
	}
	return nil
}

func (r *RemainingEntriesLogRepository) GetRemainingEntriesLogs() ([]model.RemainingEntriesLog, error) {
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

func (r *RemainingEntriesLogRepository) GetRemainingEntriesLogsByUserID(userID int) ([]model.RemainingEntriesLog, error) {
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

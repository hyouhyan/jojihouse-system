package repository

import (
	"context"
	"time"

	"jojihouse-entrance-system/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RemainingEntriesLogRepository struct {
	db *mongo.Database
}

func NewRemainingEntriesLogRepository(db *mongo.Database) *RemainingEntriesLogRepository {
	return &RemainingEntriesLogRepository{db: db}
}

func (r *RemainingEntriesLogRepository) CreateRemainingEntriesLog(log *model.RemainingEntriesLog) error {
	log.ID = primitive.NilObjectID
	log.UpdatedAt = time.Now()

	_, err := r.db.Collection("remaining_entries_log").InsertOne(context.Background(), log)
	if err != nil {
		return err
	}
	return nil
}

func (r *RemainingEntriesLogRepository) GetRemainingEntriesLogs(lastID primitive.ObjectID, limit int64) ([]model.RemainingEntriesLog, error) {
	// フィルターなしで全ログを取得
	return r._findRemainingEntriesLogs(bson.D{}, lastID, limit)
}

func (r *RemainingEntriesLogRepository) GetRemainingEntriesLogsOnlyIncrease(lastID primitive.ObjectID, limit int64) ([]model.RemainingEntriesLog, error) {
	// previous_entries より new_entries の方が大きいデータを取得
	filter := bson.D{
		{Key: "$expr", Value: bson.D{
			{Key: "$gt", Value: bson.A{"$new_entries", "$previous_entries"}},
		}},
	}

	return r._findRemainingEntriesLogs(filter, lastID, limit)
}

func (r *RemainingEntriesLogRepository) GetRemainingEntriesLogsByUserID(userID int, lastID primitive.ObjectID, limit int64) ([]model.RemainingEntriesLog, error) {
	// `user_id` でフィルター
	filter := bson.D{{Key: "user_id", Value: userID}}
	return r._findRemainingEntriesLogs(filter, lastID, limit)
}

func (r *RemainingEntriesLogRepository) GetLastRemainingEntriesLogByUserID(userID int) (model.RemainingEntriesLog, error) {
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "updated_at", Value: 1}}) // `updated_at` で昇順ソート

	filter := bson.D{{Key: "user_id", Value: userID}}

	cursor := r.db.Collection("remaining_entries_log").FindOne(context.Background(), filter, findOptions)

	var result model.RemainingEntriesLog

	err := cursor.Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.RemainingEntriesLog{}, nil
		}
		return model.RemainingEntriesLog{}, err
	}

	return result, nil
}

// 共通の検索処理
func (r *RemainingEntriesLogRepository) _findRemainingEntriesLogs(filter bson.D, lastID primitive.ObjectID, limit int64) ([]model.RemainingEntriesLog, error) {
	var logs []model.RemainingEntriesLog

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{Key: "updated_at", Value: 1}}) // `updated_at` で昇順ソート

	// lastID によるページネーション
	if !lastID.IsZero() {
		filter = append(filter, bson.E{Key: "_id", Value: bson.D{{Key: "$gt", Value: lastID}}})
	}

	cursor, err := r.db.Collection("remaining_entries_log").Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var log model.RemainingEntriesLog
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}

		// タイムゾーンの変換
		log.UpdatedAt = log.UpdatedAt.In(time.Local)

		logs = append(logs, log)
	}
	return logs, nil
}

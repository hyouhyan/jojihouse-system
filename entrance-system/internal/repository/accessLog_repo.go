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

func (r *AccessLogRepository) GetAccessLogs(lastID primitive.ObjectID) ([]model.AccessLog, error) {
	// フィルターなしで全ログを取得
	return r._findAccessLogs(bson.D{}, lastID, 50)
}

func (r *AccessLogRepository) GetAccessLogsByUserID(userID int, lastID primitive.ObjectID) ([]model.AccessLog, error) {
	// `user_id` でフィルター
	filter := bson.D{{Key: "user_id", Value: userID}}
	return r._findAccessLogs(filter, lastID, 50)
}

// 共通の検索処理
func (r *AccessLogRepository) _findAccessLogs(filter bson.D, lastID primitive.ObjectID, limit int64) ([]model.AccessLog, error) {
	var logs []model.AccessLog

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{Key: "time", Value: 1}}) // `time` で昇順ソート

	// lastID によるページネーション
	if !lastID.IsZero() {
		filter = append(filter, bson.E{Key: "_id", Value: bson.D{{Key: "$gt", Value: lastID}}})
	}

	cursor, err := r.db.Collection("access_log").Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var log model.AccessLog
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}

		// タイムゾーンの変換
		log.Time = log.Time.In(time.Local)

		logs = append(logs, log)
	}
	return logs, nil
}

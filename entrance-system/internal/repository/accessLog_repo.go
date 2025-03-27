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
	var logs []model.AccessLog

	findOptions := options.Find()
	findOptions.SetLimit(50)
	findOptions.SetSort(bson.D{{Key: "time", Value: 1}}) // `time` で昇順ソート

	// lastID がゼロ値でなければ、それ以降のデータを取得
	filter := bson.D{}
	if !lastID.IsZero() {
		filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: lastID}}}} // `_id` を基準にページネーション
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

func (r *AccessLogRepository) GetAccessLogsByUserID(userID int, lastID primitive.ObjectID) ([]model.AccessLog, error) {
	var logs []model.AccessLog

	findOptions := options.Find()
	findOptions.SetLimit(50)
	findOptions.SetSort(bson.D{{Key: "time", Value: 1}}) // `time` で昇順ソート

	// lastID がゼロ値でなければ、それ以降のデータを取得
	filter := bson.D{{Key: "user_id", Value: userID}}
	if !lastID.IsZero() {
		filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: lastID}}}} // `_id` を基準にページネーション
	}

	cursor, err := r.db.Collection("access_log").Find(context.Background(), filter)
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

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

func (r *RemainingEntriesLogRepository) CreateFixedRemainingEntriesLog(log *model.RemainingEntriesLog) error {
	log.ID = primitive.NilObjectID

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
	findOptions.SetSort(bson.D{{Key: "updated_at", Value: -1}}) // これは-1が正しい、mongo compassで確認してみて

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

	// lastIDからデータを取得して、lastTimeを設定
	var lastTime time.Time
	if !lastID.IsZero() {
		// lastIDから取得できるタイムスタンプとupdated_atフィールドの値は異なる可能性がある
		// よって、ドキュメントから取得する必要がある
		var lastLog model.RemainingEntriesLog
		err := r.db.Collection("remaining_entries_log").FindOne(context.Background(), bson.D{{Key: "_id", Value: lastID}}).Decode(&lastLog)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				lastTime = time.Time{} // lastIDが存在しない場合はゼロ値を設定
				filter = append(filter, bson.E{Key: "_id", Value: bson.D{{Key: "$lt", Value: lastID}}})
			} else {
				return nil, err // 他のエラーはそのまま返す
			}
		} else {
			lastTime = lastLog.UpdatedAt.In(time.Local) // lastIDから時間を取得し、ローカルタイムゾーンに変換
		}
	} else {
		// lastIDがゼロの場合は、lastTimeをゼロに設定
		lastTime = time.Time{}
	}

	// lastTime によるページネーション
	if !lastTime.IsZero() {
		filter = append(filter, bson.E{Key: "updated_at", Value: bson.D{{Key: "$lt", Value: lastTime}}})
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

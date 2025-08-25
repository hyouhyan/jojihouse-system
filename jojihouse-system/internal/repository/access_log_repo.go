package repository

import (
	"context"
	"time"

	"jojihouse-system/internal/model"

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

func (r *AccessLogRepository) CreateEntryAccessLog(userid int) error {
	log :=
		&model.AccessLog{
			UserID:     userid,
			Time:       time.Now(),
			AccessType: "entry",
		}

	return r.CreateAccessLog(log)
}

func (r *AccessLogRepository) CreateExitAccessLog(userid int) error {
	log :=
		&model.AccessLog{
			UserID:     userid,
			Time:       time.Now(),
			AccessType: "exit",
		}

	return r.CreateAccessLog(log)
}

func (r *AccessLogRepository) GetAccessLogsByAnyFilter(lastID primitive.ObjectID, options model.AccessLogFilter) ([]model.AccessLog, error) {
	filter := bson.D{}

	// UserID のフィルタ
	if options.UserID > 0 {
		filter = append(filter, bson.E{Key: "user_id", Value: options.UserID})
	}

	// 日時フィルタ
	timeFilter := bson.D{}
	if !options.DayAfter.IsZero() {
		timeFilter = append(timeFilter, bson.E{Key: "$gte", Value: options.DayAfter})
	}
	if !options.DayBefore.IsZero() {
		timeFilter = append(timeFilter, bson.E{Key: "$lte", Value: options.DayBefore})
	}
	if len(timeFilter) > 0 {
		filter = append(filter, bson.E{Key: "time", Value: timeFilter})
	}

	// AccessType のフィルタ
	if options.AccessType != "" {
		filter = append(filter, bson.E{Key: "access_type", Value: options.AccessType})
	}

	// リミットの設定（デフォルト50）
	limit := int64(50)
	if options.Limit > 0 {
		limit = int64(options.Limit)
	}

	return r._findAccessLogs(filter, lastID, limit)
}

// 最終アクセスログを取得
func (r *AccessLogRepository) GetLastAccessLogByUserID(userID int) (*model.AccessLog, error) {
	filter := bson.D{
		{Key: "user_id", Value: userID},
	}

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}}) // `time` で降順ソート

	var log model.AccessLog
	err := r.db.Collection("access_log").FindOne(context.Background(), filter, findOptions).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // ドキュメントが見つからない場合は nil を返す
		}
		return nil, err
	}

	log.Time = log.Time.In(time.Local) // タイムゾーンの変換
	return &log, nil
}

// 共通の検索処理
func (r *AccessLogRepository) _findAccessLogs(filter bson.D, lastID primitive.ObjectID, limit int64) ([]model.AccessLog, error) {
	var logs []model.AccessLog

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{Key: "time", Value: -1}}) // `time` で昇順ソート

	// lastIDからデータを取得して、timeフィールドからlastTimeを設定
	var lastTime time.Time
	if !lastID.IsZero() {
		// lastIDから取得できるタイムスタンプとtimeフィールドの値は異なる可能性がある
		// よって、ドキュメントから取得する必要がある
		var lastLog model.AccessLog
		err := r.db.Collection("access_log").FindOne(context.Background(), bson.D{{Key: "_id", Value: lastID}}).Decode(&lastLog)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				lastTime = time.Time{} // lastIDが存在しない場合はゼロ値を設定
			} else {
				return nil, err // 他のエラーはそのまま返す
			}
		} else {
			lastTime = lastLog.Time.In(time.Local) // lastIDから時間を取得し、ローカルタイムゾーンに変換
		}
	} else {
		// lastIDがゼロの場合は、lastTimeをゼロに設定
		lastTime = time.Time{}
	}

	// lastTime によるページネーション
	if !lastTime.IsZero() {
		filter = append(filter, bson.E{Key: "time", Value: bson.D{{Key: "$lt", Value: lastTime}}})
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

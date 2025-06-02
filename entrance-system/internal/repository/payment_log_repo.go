package repository

import (
	"context"
	"fmt"
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
	opts.SetSort(bson.D{{Key: "time", Value: -1}})
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

func (r *PaymentLogRepository) GetMonthlyPaymentLogs(year int, month int) ([]model.PaymentLog, int, error) {
	// contextの準備 (関数全体で利用できるように最初に定義)
	ctx := context.Background() // 必要に応じてタイムアウト付きコンテキストなどに置き換え

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	// 月末の最終日時を正確に設定するため、翌月の初日を取得し、それより前とする
	endDate := startDate.AddDate(0, 1, 0) // 翌月の1日の0時0分0秒

	filter := bson.D{
		{Key: "time", Value: bson.D{
			{Key: "$gte", Value: startDate}, // 指定月の1日 00:00:00 以降
			{Key: "$lt", Value: endDate},    // 翌月の1日 00:00:00 より前
		}},
	}

	// 合計値の計算 (MongoDBサーバーサイド)
	var totalAmount int
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}}, // 期間でフィルタ
		bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: nil}, // 全ドキュメントを1グループに
				{Key: "totalAmount", Value: bson.D{{Key: "$sum", Value: "$amount"}}}, // amountの合計
			},
		}},
	}

	cursorTotal, err := r.db.Collection("payment_log").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("合計値の集計クエリ実行に失敗: %w", err)
	}
	defer cursorTotal.Close(ctx)

	var aggResults []totalAmountResult
	if err = cursorTotal.All(ctx, &aggResults); err != nil {
		return nil, 0, fmt.Errorf("合計値の集計結果デコードに失敗: %w", err)
	}

	if len(aggResults) > 0 {
		totalAmount = aggResults[0].TotalAmount
	}
	// aggResultsが空の場合 (該当期間にドキュメントがない場合)、totalAmountは0のまま

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "time", Value: -1}}) // 時間の降順でソート

	cursor, err := r.db.Collection("payment_log").Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err // エラー発生時はログリストnil, total 0, エラーを返す
	}
	defer cursor.Close(ctx)

	var logs []model.PaymentLog

	for cursor.Next(ctx) {
		var logEntry model.PaymentLog // デコード先の変数
		if err := cursor.Decode(&logEntry); err != nil {
			// デコードエラーが発生した場合、それまでのログと合計は返さずエラーを返す
			return nil, 0, err
		}
		logs = append(logs, logEntry)
	}

	if err := cursor.Err(); err != nil {
		// カーソル処理中にエラーが発生した場合
		return nil, 0, err
	}

	// 取得したログのスライスと、計算したamountの合計値を返す
	return logs, totalAmount, nil
}

// 集計結果（合計値）をマッピングするための構造体
type totalAmountResult struct {
	ID          interface{} `bson:"_id"` // グループ化キー (nilが入る)
	TotalAmount int         `bson:"totalAmount"`
}

package repository

import (
	"context"
	"fmt"
	"jojihouse-system/internal/model"
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

// MonthlyTotalAmount は月次の支払い合計額を表す構造体です
type MonthlyTotalAmount struct {
	Total      int
	OliveTotal int
	CashTotal  int
}

// getMonthlyTotalAmount は指定された年月の支払い合計額を取得します
func (r *PaymentLogRepository) getMonthlyTotalAmount(year int, month int) (*MonthlyTotalAmount, error) {
	ctx := context.Background()
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0)

	filter := bson.D{
		{Key: "time", Value: bson.D{
			{Key: "$gte", Value: startDate},
			{Key: "$lt", Value: endDate},
		}},
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: "$payway"},
				{Key: "total", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
			},
		}},
	}

	cursorTotal, err := r.db.Collection("payment_log").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("合計値の集計クエリ実行に失敗: %w", err)
	}
	defer cursorTotal.Close(ctx)

	var aggResults []struct {
		ID    string `bson:"_id"`
		Total int    `bson:"total"`
	}
	if err = cursorTotal.All(ctx, &aggResults); err != nil {
		return nil, fmt.Errorf("合計値の集計結果デコードに失敗: %w", err)
	}

	totals := &MonthlyTotalAmount{}
	for _, aggResult := range aggResults {
		totals.Total += aggResult.Total
		switch aggResult.ID {
		case "olive":
			totals.OliveTotal = aggResult.Total
		case "cash":
			totals.CashTotal = aggResult.Total
		}
	}

	return totals, nil
}

func (r *PaymentLogRepository) GetMonthlyPaymentLogs(year int, month int) (*model.MonthlyPaymentLog, error) {
	ctx := context.Background()
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0)

	filter := bson.D{
		{Key: "time", Value: bson.D{
			{Key: "$gte", Value: startDate},
			{Key: "$lt", Value: endDate},
		}},
	}

	totals, err := r.getMonthlyTotalAmount(year, month)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "time", Value: -1}})

	cursor, err := r.db.Collection("payment_log").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []model.PaymentLog
	for cursor.Next(ctx) {
		var logEntry model.PaymentLog
		if err := cursor.Decode(&logEntry); err != nil {
			return nil, err
		}
		logs = append(logs, logEntry)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &model.MonthlyPaymentLog{
		Year:       year,
		Month:      month,
		Total:      totals.Total,
		OliveTotal: totals.OliveTotal,
		CashTotal:  totals.CashTotal,
		Logs:       logs,
	}, nil
}

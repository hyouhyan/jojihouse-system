package main

import (
	"context"
	"fmt"
	"jojihouse-system-consistency-checker/database"
	"jojihouse-system-consistency-checker/model"
	"time"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// DBへ接続
	database.ConnectPostgres()
	defer database.ClosePostgres()

	database.ConnectMongo()
	defer database.CloseMongo()

	users, err := GetAllUsers(database.PostgresDB)
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		return
	}

	for _, user := range users {
		year := 2025
		month := 8
		logs, err := GetUsersAccessLog(database.MongoDB, *user.ID, year, month)
		if err != nil {
			fmt.Printf("Error fetching access logs for user %d: %v\n", user.Number, err)
			continue
		}

		for _, log := range logs {
			fmt.Printf("User %s accessed at %s\n", *user.Name, log.Time)
		}
	}
}

func GetAllUsers(db *sqlx.DB) ([]model.User, error) {
	var users []model.User
	err := db.Select(&users, "SELECT * FROM users ORDER BY number")

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func GetUsersAccessLog(db *mongo.Database, userID int, year int, month int) ([]model.AccessLog, error) {
	collection := db.Collection("access_log")
	filter := map[string]interface{}{
		"user_id": userID,
		"time": map[string]interface{}{
			"$gte": time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local),
			"$lt":  time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.Local),
		},
	}

	options := options.Find().SetSort(map[string]int{"time": 1})

	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, fmt.Errorf("failed to find access logs: %w", err)
	}
	defer cursor.Close(context.Background())

	var logs []model.AccessLog
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

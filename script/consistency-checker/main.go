package main

import (
	"context"
	"fmt"
	"jojihouse-system-consistency-checker/database"
	"jojihouse-system-consistency-checker/model"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type usersAccessCount struct {
	username      string
	daysPassed    int
	accesCount    int
	decreaseCount int
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <year> <month>")
		return
	}

	yearStr := os.Args[1]
	monthStr := os.Args[2]

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Printf("Invalid year: %v\n", err)
		return
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		log.Printf("Invalid month: %v\n", err)
		return
	}

	// DBへ接続
	database.ConnectPostgres()
	defer database.ClosePostgres()

	database.ConnectMongo()
	defer database.CloseMongo()

	users, err := GetAllUsers(database.PostgresDB)
	if err != nil {
		log.Printf("Error fetching users: %v\n", err)
		return
	}

	var accessCounts []usersAccessCount

	for _, user := range users {
		logs, err := GetUsersAccessLog(database.MongoDB, *user.ID, year, month)
		if err != nil {
			log.Printf("Error fetching access logs for user %d: %v\n", user.Number, err)
			continue
		}

		var lastEntryLog *model.AccessLog
		accessCount := 0
		daysPassed := 0
		for _, log := range logs {
			if log.AccessType == "entry" {
				if lastEntryLog == nil {
					accessCount++
				} else {
					if !isSameDate(log.Time, lastEntryLog.Time) {
						accessCount++
					}
				}
				lastEntryLog = &log
			}
			if log.AccessType == "exit" {
				if lastEntryLog != nil {
					if !isSameDate(log.Time, lastEntryLog.Time) {
						accessCount += getPassedDays(lastEntryLog.Time, log.Time)
						daysPassed += getPassedDays(lastEntryLog.Time, log.Time)
						lastEntryLog = &log
					}
				}
			}
		}
		count, err := GetRemainingEntriesLogsDecreaseCount(database.MongoDB, *user.ID, year, month)
		if err != nil {
			log.Printf("Error fetching remaining entries logs: %v\n", err)
			return
		}
		accessCounts = append(accessCounts, usersAccessCount{
			username:      *user.Name,
			daysPassed:    daysPassed,
			accesCount:    accessCount,
			decreaseCount: count,
		})
	}

	// CSV出力
	fmt.Println("ユーザー名,経過日数,入場回数,入場可能回数減少分")
	for _, ac := range accessCounts {
		fmt.Printf("%s,%d,%d,%d\n", ac.username, ac.daysPassed, ac.accesCount, ac.decreaseCount)
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

func isSameDate(a, b time.Time) bool {
	// bのtimezoneをaのtimezoneに合わせる
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}

	aDate := cnvTo00Time(a)
	bDate := cnvTo00Time(b)

	return aDate.Equal(bDate)
}

func cnvTo00Time(t time.Time) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0, t.Location())
}

func getPassedDays(targetDate, currentDate time.Time) int {
	// aとbのtimezoneを揃える
	if targetDate.Location() != currentDate.Location() {
		currentDate = currentDate.In(targetDate.Location())
	}

	// 00:00:00どうしで比較
	targetDate = cnvTo00Time(targetDate)
	currentDate = cnvTo00Time(currentDate)

	// 日数の差を計算
	daysPassed := int(currentDate.Sub(targetDate).Hours() / 24)

	return daysPassed
}

func GetRemainingEntriesLogsDecreaseCount(db *mongo.Database, userID int, year int, month int) (int, error) {
	collection := db.Collection("remaining_entries_log")
	filter := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "$expr", Value: bson.D{
			{Key: "$gt", Value: bson.A{"$previous_entries", "$new_entries"}},
		}},
		{Key: "updated_at", Value: bson.D{
			{Key: "$gte", Value: time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)},
			{Key: "$lt", Value: time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.Local)},
		}},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to find access logs: %w", err)
	}
	defer cursor.Close(context.Background())

	count := 0
	for cursor.Next(context.Background()) {
		var log model.RemainingEntriesLog
		if err := cursor.Decode(&log); err != nil {
			return 0, err
		}

		count += log.NewEntries - log.PreviousEntries
	}

	count = int(math.Abs(float64(count)))

	return count, nil
}

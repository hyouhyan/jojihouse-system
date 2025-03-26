package main

import (
	"database/sql"
	"log"
	"time"

	"jojihouse-entrance-system/internal/database"
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
	"jojihouse-entrance-system/internal/service"
)

func main() {
	database.Connect()
	defer database.Close()

	test(database.DB)
}

func test(db *sql.DB) {
	// ユーザーリポジトリを作成
	userRepo := repository.NewUserRepository(db)

	// ユーザーサービスを作成
	userService := service.NewUserService(userRepo)

	// ユーザーを追加
	// ユーザー情報を作成
	user := &model.User{
		Name:              "test",
		Description:       "テストユーザー",
		Barcode:           "test",
		Contact:           "test",
		Remaining_entries: 10,
		Registered_at:     time.Now(),
		Total_entries:     0,
	}

	// ユーザーを追加
	_, err := userService.CreateUser(user)
	if err != nil {
		log.Fatalf("Failed to create a user: %v", err)
	}
}

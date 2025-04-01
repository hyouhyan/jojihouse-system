package main

import (
	"jojihouse-entrance-system/api/handler"
	"jojihouse-entrance-system/api/router"
	"jojihouse-entrance-system/internal/database"
	"jojihouse-entrance-system/internal/repository"
	"jojihouse-entrance-system/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectPostgres()
	defer database.ClosePostgres()

	database.ConnectMongo()
	defer database.CloseMongo()

	// 依存関係の注入

	// ユーザーリポジトリを作成
	userRepo := repository.NewUserRepository(database.PostgresDB)
	// ロールリポジトリを作成
	roleRepo := repository.NewRoleRepository(database.PostgresDB)
	// ログリポジトリを作成
	accessLogRepo := repository.NewLogRepository(database.MongoDB)
	// 入場可能回数ログリポジトリ
	remainingEntriesLogRepo := repository.NewRemainingEntriesLogRepository(database.MongoDB)
	// 在室ユーザーリポジトリ
	currentUsersRepo := repository.NewCurrentUsersRepository(database.PostgresDB)

	// entranceサービスを作成
	entranceService := service.NewEntranceService(userRepo, roleRepo, accessLogRepo, remainingEntriesLogRepo, currentUsersRepo)
	// adminManagementサービスを作成
	// adminManagementService := service.NewAdminManagementService(userRepo, roleRepo, accessLogRepo, remainingEntriesLogRepo)
	// userPortalサービスを作成
	// userPortalService := service.NewUserPortalService(userRepo, roleRepo, accessLogRepo, remainingEntriesLogRepo, currentUsersRepo)

	// EntranceHandlerを作成
	entranceHandler := handler.NewEntranceHandler(entranceService)

	// Gin ルーターの設定
	r := gin.Default()
	router.SetupEntranceRoutes(r, entranceHandler)

	// test(database.PostgresDB, database.MongoDB)

	// サーバー起動
	r.Run("127.0.0.1:8080")
}

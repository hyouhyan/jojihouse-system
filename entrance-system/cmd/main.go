package main

import (
	"jojihouse-entrance-system/api/handler"
	"jojihouse-entrance-system/api/router"
	"jojihouse-entrance-system/internal/database"
	"jojihouse-entrance-system/internal/repository"
	"jojihouse-entrance-system/internal/service"
	"time"

	_ "jojihouse-entrance-system/docs" // 生成される Swagger ドキュメントを読み込む

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title JojiHouse Entrance System API
// @version 1.0
// @description JojiHouse の入退室管理システム API ドキュメント
// @host 127.0.0.1:8080
// @BasePath /
func main() {
	database.ConnectPostgres()
	defer database.ClosePostgres()

	database.ConnectMongo()
	defer database.CloseMongo()

	// 依存関係の注入

	// ユーザーリポジトリ
	userRepo := repository.NewUserRepository(database.PostgresDB)
	// ロールリポジトリ
	roleRepo := repository.NewRoleRepository(database.PostgresDB)
	// ログリポジトリ
	accessLogRepo := repository.NewLogRepository(database.MongoDB)
	// 入場可能回数ログリポジトリ
	remainingEntriesLogRepo := repository.NewRemainingEntriesLogRepository(database.MongoDB)
	// 在室ユーザーリポジトリ
	currentUsersRepo := repository.NewCurrentUsersRepository(database.PostgresDB)

	// entranceサービス
	entranceService := service.NewEntranceService(userRepo, roleRepo, accessLogRepo, remainingEntriesLogRepo, currentUsersRepo)
	// adminManagementサービス
	adminManagementService := service.NewAdminManagementService(userRepo, roleRepo, accessLogRepo, remainingEntriesLogRepo)
	// userPortalサービス
	userPortalService := service.NewUserPortalService(userRepo, roleRepo, accessLogRepo, remainingEntriesLogRepo, currentUsersRepo)

	// EntranceHandler
	entranceHandler := handler.NewEntranceHandler(entranceService, userPortalService)
	userHandler := handler.NewUserHandler(userPortalService, adminManagementService)
	roleHandler := handler.NewRoleHandler(userPortalService)
	kaisukenHandler := handler.NewKaisukenHandler(userPortalService, adminManagementService)

	// Gin ルーターの設定
	r := gin.Default()

	// CORSの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://127.0.0.1:8080",
			"http://house.joji:8080",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // プリフライトリクエストの結果をキャッシュ
	}))

	router.SetupEntranceRoutes(r, entranceHandler)
	router.SetupUserRoutes(r, userHandler)
	router.SetupRoleRoutes(r, roleHandler)
	router.SetupKaisukenRoutes(r, kaisukenHandler)

	// test(database.PostgresDB, database.MongoDB)

	// Swagger のエンドポイントを追加
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// サーバー起動
	r.Run("0.0.0.0:8080")
}

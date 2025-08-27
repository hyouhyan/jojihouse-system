package main

import (
	"jojihouse-system/api/authentication"
	"jojihouse-system/api/handler"
	"jojihouse-system/api/router"
	"jojihouse-system/internal/database"
	"jojihouse-system/internal/repository"
	"jojihouse-system/internal/service"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// DBへ接続
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
	// 支払いログリポジトリ
	paymentLogRepo := repository.NewPaymentLogRepository(database.MongoDB)
	// Discord通知
	discordNoticeRepository := repository.NewDiscordNoticeRepository()

	// entranceサービス
	entranceService := service.NewEntranceService(userRepo, roleRepo, accessLogRepo, remainingEntriesLogRepo, currentUsersRepo, discordNoticeRepository)
	// adminManagementサービス
	adminManagementService := service.NewAdminManagementService(userRepo, roleRepo, accessLogRepo, remainingEntriesLogRepo, paymentLogRepo)
	// userPortalサービス
	userPortalService := service.NewUserPortalService(userRepo, roleRepo, accessLogRepo, remainingEntriesLogRepo, currentUsersRepo)

	// Discord Authentication
	// discordAuthentication := authentication.NewDiscordAuthentication(userPortalService)
	discordAuthentication := authentication.NewDiscordAuthentication()

	// EntranceHandler
	entranceHandler := handler.NewEntranceHandler(entranceService, userPortalService)
	userHandler := handler.NewUserHandler(userPortalService, adminManagementService)
	roleHandler := handler.NewRoleHandler(userPortalService)
	kaisukenHandler := handler.NewKaisukenHandler(userPortalService, adminManagementService)
	paymentHandler := handler.NewPaymentHandler(adminManagementService)
	authHandler := handler.NewAuthHandler(discordAuthentication)

	// Gin ルーターの設定
	r := gin.Default()

	// 環境変数から許可するオリジンのリストを取得
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	// カンマ区切りの文字列を文字列スライスに変換
	originsList := strings.Split(allowedOrigins, ",")

	// CORSの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     originsList,
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
	router.SetupPaymentRoutes(r, paymentHandler)
	router.SetupAuthRoutes(r, authHandler)

	// サーバー起動
	r.Run("0.0.0.0:8080")
}

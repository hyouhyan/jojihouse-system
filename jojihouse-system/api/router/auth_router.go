package router

import (
	"jojihouse-system/api/authentication/middleware"
	"jojihouse-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, handler *handler.AuthHandler, middleware *middleware.AuthMiddleware) {
	authGroup := router.Group("/auth")
	{
		authGroup.GET("discord", handler.DiscordAuth)
	}
}

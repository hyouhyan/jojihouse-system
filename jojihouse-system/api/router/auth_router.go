package router

import (
	"jojihouse-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, handler *handler.AuthHandler) {
	authGroup := router.Group("/auth")
	{
		authGroup.GET("discord", handler.DiscordAuth)
	}
}

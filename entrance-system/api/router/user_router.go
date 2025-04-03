package router

import (
	"jojihouse-entrance-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.GET("/", userHandler.GetAllUsers)
	}
}

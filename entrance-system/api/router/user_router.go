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
		userGroup.GET("/:user_id", userHandler.GetUserByID)
		userGroup.PUT("/:user_id", userHandler.UpdateUser)
		userGroup.DELETE("/:user_id", userHandler.DeleteUser)
		userGroup.GET("/:user_id/roles", userHandler.GetRolesByUserID)
		userGroup.POST("/:user_id/roles", userHandler.AddRoleToUser)
		userGroup.DELETE("/:user_id/roles/:role_id", userHandler.RemoveRoleFromUser)
	}
}

package router

import (
	"jojihouse-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userHandler *handler.UserHandler) {
	userGroupMember := router.Group("/users")
	// userGroupMember.Use()
	{
		userGroupMember.GET("", userHandler.GetUsers)
		userGroupMember.GET("/:user_id", userHandler.GetUserByID)

		userGroupMember.GET("/:user_id/roles", userHandler.GetRolesByUserID)
	}

	userGroupAdmin := router.Group("/users")
	// userGroupAdmin.Use()
	{
		userGroupAdmin.POST("", userHandler.CreateUser)

		userGroupAdmin.PATCH("/:user_id", userHandler.UpdateUser)
		userGroupAdmin.DELETE("/:user_id", userHandler.DeleteUser)

		userGroupAdmin.POST("/:user_id/roles", userHandler.AddRoleToUser)
		userGroupAdmin.DELETE("/:user_id/roles/:role_id", userHandler.RemoveRoleFromUser)

		userGroupAdmin.GET("/:user_id/logs", userHandler.GetUserLogs)
	}
}

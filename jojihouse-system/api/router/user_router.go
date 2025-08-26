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

	userGroupHouseAdmin := router.Group("/users")
	// userGroupHouseAdmin.Use()
	{
		userGroupHouseAdmin.PATCH("/:user_id", userHandler.UpdateUser)

		userGroupHouseAdmin.GET("/:user_id/logs", userHandler.GetUserLogs)
	}

	userGroupSysAdmin := router.Group("/users")
	// userGroupSysAdmin.Use()
	{
		userGroupSysAdmin.POST("", userHandler.CreateUser)

		userGroupSysAdmin.DELETE("/:user_id", userHandler.DeleteUser)

		userGroupSysAdmin.POST("/:user_id/roles", userHandler.AddRoleToUser)
		userGroupSysAdmin.DELETE("/:user_id/roles/:role_id", userHandler.RemoveRoleFromUser)
	}
}

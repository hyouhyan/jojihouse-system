package router

import (
	"jojihouse-system/api/authentication/middleware"
	"jojihouse-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userHandler *handler.UserHandler, middleware *middleware.AuthMiddleware) {
	userGroupMember := router.Group("/users")
	userGroupMember.Use(middleware.AuthMember)
	{
		userGroupMember.GET("", userHandler.GetUsers)
		userGroupMember.GET("/:user_id", userHandler.GetUserByID)

		userGroupMember.GET("/:user_id/roles", userHandler.GetRolesByUserID)
	}

	userGroupHouseAdmin := router.Group("/users")
	userGroupHouseAdmin.Use(middleware.AuthHouseAdmin)
	{
		userGroupHouseAdmin.PATCH("/:user_id", userHandler.UpdateUser)

		userGroupHouseAdmin.GET("/:user_id/logs", userHandler.GetUserLogs)
	}

	userGroupSysAdmin := router.Group("/users")
	userGroupSysAdmin.Use(middleware.AuthSystemAdmin)
	{
		userGroupSysAdmin.POST("", userHandler.CreateUser)

		userGroupSysAdmin.DELETE("/:user_id", userHandler.DeleteUser)

		userGroupSysAdmin.POST("/:user_id/roles", userHandler.AddRoleToUser)
		userGroupSysAdmin.DELETE("/:user_id/roles/:role_id", userHandler.RemoveRoleFromUser)
	}
}

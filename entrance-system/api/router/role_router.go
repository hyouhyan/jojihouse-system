package router

import (
	"jojihouse-entrance-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoleRoutes(router *gin.Engine, roleHandler *handler.RoleHandler) {
	userGroup := router.Group("/roles")
	{
		userGroup.GET("/", roleHandler.GetAllRoles)
	}
}

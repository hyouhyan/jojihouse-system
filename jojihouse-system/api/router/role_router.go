package router

import (
	"jojihouse-system/api/authentication/middleware"
	"jojihouse-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoleRoutes(router *gin.Engine, roleHandler *handler.RoleHandler, middleware *middleware.AuthMiddleware) {
	roleGroupMember := router.Group("/roles")
	// roleGroupMember.Use()
	{
		roleGroupMember.GET("", roleHandler.GetAllRoles)
	}
}

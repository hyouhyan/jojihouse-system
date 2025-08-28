package router

import (
	"jojihouse-system/api/authentication/middleware"
	"jojihouse-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupKaisukenRoutes(router *gin.Engine, handler *handler.KaisukenHandler, middleware *middleware.AuthMiddleware) {
	kaisukenGroupHouseAdmin := router.Group("/kaisuken")
	kaisukenGroupHouseAdmin.Use(middleware.AuthHouseAdmin)
	{
		kaisukenGroupHouseAdmin.POST("", handler.BuyKaisuken)
	}
}

package router

import (
	"jojihouse-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupKaisukenRoutes(router *gin.Engine, handler *handler.KaisukenHandler) {
	kaisukenGroupHouseAdmin := router.Group("/kaisuken")
	// kaisukenGroupHouseAdmin.Use()
	{
		kaisukenGroupHouseAdmin.POST("", handler.BuyKaisuken)
	}
}

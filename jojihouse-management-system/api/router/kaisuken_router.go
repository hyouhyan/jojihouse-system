package router

import (
	"jojihouse-management-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupKaisukenRoutes(router *gin.Engine, handler *handler.KaisukenHandler) {
	group := router.Group("/kaisuken")
	{
		group.POST("", handler.BuyKaisuken)
	}
}

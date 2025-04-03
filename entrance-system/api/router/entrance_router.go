package router

import (
	"jojihouse-entrance-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupEntranceRoutes(router *gin.Engine, entranceHandler *handler.EntranceHandler) {
	entranceGroup := router.Group("/entrance")
	{
		entranceGroup.POST("", entranceHandler.RecordEntrance)
		entranceGroup.GET("/current", entranceHandler.GetCurrentUsers)
		entranceGroup.GET("/logs", entranceHandler.GetAccessLogs)
		entranceGroup.GET("/logs/:user_id", entranceHandler.GetAccessLogsByUserID)
	}
}

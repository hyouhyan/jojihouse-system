package router

import (
	"jojihouse-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupEntranceRoutes(router *gin.Engine, entranceHandler *handler.EntranceHandler) {
	entranceGroupMember := router.Group("/entrance")
	// entranceGroupMember.Use()
	{
		entranceGroupMember.GET("/current", entranceHandler.GetCurrentUsers)
	}

	entranceGroupHouseAdmin := router.Group("/entrance")
	// entranceGroupHouseAdmin.Use()
	{
		entranceGroupHouseAdmin.GET("/logs", entranceHandler.GetAccessLogs)
		entranceGroupHouseAdmin.GET("/logs/:user_id", entranceHandler.GetAccessLogsByUserID)
	}

	entranceGroupEntrance := router.Group("/entrance")
	// entranceGroupEntrance.Use()
	{
		entranceGroupEntrance.POST("", entranceHandler.RecordEntrance)
	}
}

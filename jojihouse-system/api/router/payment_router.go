package router

import (
	"jojihouse-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(router *gin.Engine, handler *handler.PaymentHandler) {
	group := router.Group("/payment")
	{
		group.GET("", handler.GetAllPaymentLogs)
		group.POST("", handler.CreatePaymentLog)

		group.GET("/:log_id", handler.GetPaymentLogByID)

		group.GET("/monthly", handler.GetMonthlyPaymentLogs)
	}
}

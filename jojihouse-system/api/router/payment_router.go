package router

import (
	"jojihouse-system/api/authentication/middleware"
	"jojihouse-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(router *gin.Engine, handler *handler.PaymentHandler, middleware *middleware.AuthMiddleware) {
	paymentGroupHouseAdmin := router.Group("/payment")
	// paymentGroupHouseAdmin.Use()
	{
		paymentGroupHouseAdmin.GET("", handler.GetAllPaymentLogs)
		paymentGroupHouseAdmin.POST("", handler.CreatePaymentLog)
		paymentGroupHouseAdmin.GET("/monthly", handler.GetMonthlyPaymentLogs)
	}
}

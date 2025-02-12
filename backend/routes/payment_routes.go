package routes

import (
	"github.com/TgkCapture/alumni-welfare/controllers"
	"github.com/TgkCapture/alumni-welfare/middleware"
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(r *gin.Engine) {
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/pay", controllers.MakePayment)
		auth.GET("/payments", controllers.GetPaymentHistory)
		auth.POST("/webhook", controllers.PaymentWebhook)
		auth.GET("/payments/:chargeId/verify", controllers.VerifyPayment)
	}
}

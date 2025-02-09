package routes

import (
	"github.com/TgkCapture/alumni-welfare/controllers"
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(r *gin.Engine) {
	r.POST("/pay", controllers.MakePayment)
}

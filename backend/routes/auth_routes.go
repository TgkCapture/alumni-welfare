package routes

import (
	"github.com/TgkCapture/alumni-welfare/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Alumni Welfare Payment System API",
		})
	})
}

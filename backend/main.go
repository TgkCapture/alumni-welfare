package main

import (
	"fmt"
	"log"

	"github.com/TgkCapture/alumni-welfare/config"
	"github.com/TgkCapture/alumni-welfare/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Connect to the database
	config.ConnectDatabase()

	// Auto Migrate
	config.DB.AutoMigrate(&models.User{})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Alumni Welfare Payment System API",
		})
	})

	// Run server
	port := "8080"
	fmt.Println("Server is running on port " + port)
	log.Fatal(r.Run(":" + port))
}

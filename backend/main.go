package main

import (
	"fmt"
	"log"

	"github.com/TgkCapture/alumni-welfare/config"
	"github.com/TgkCapture/alumni-welfare/models"
	"github.com/TgkCapture/alumni-welfare/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config.SetupCORS(r)

	// Connect to the database
	config.ConnectDatabase()

	// Auto Migrate
	config.DB.AutoMigrate(&models.User{}, &models.Payment{})

	// Register routes
	routes.PaymentRoutes(r)
	routes.AuthRoutes(r)

	// Run server
	port := "8080"
	fmt.Println("Server is running on port " + port)
	log.Fatal(r.Run(":" + port))
}

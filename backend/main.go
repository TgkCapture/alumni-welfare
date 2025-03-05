package main

import (
	"fmt"
	"log"
	"os"

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
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Server is running on " + host + ":" + port)
	log.Fatal(r.Run(host + ":" + port))
}

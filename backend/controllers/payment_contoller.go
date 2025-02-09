package controllers

import (
	"net/http"

	"github.com/TgkCapture/alumni-welfare/config"
	"github.com/TgkCapture/alumni-welfare/models"
	"github.com/gin-gonic/gin"
)

type PaymentRequest struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
	Month  int    `json:"month"`
}

func MakePayment(c *gin.Context) {
	var request PaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	payment := models.Payment{
		UserID: user.ID,
		Name:   request.Name,
		Amount: request.Amount,
		Month:  request.Month,
	}

	config.DB.Create(&payment)

	c.JSON(http.StatusOK, gin.H{"message": "Payment recorded successfully!"})
}

// Get user's payment history
func GetPaymentHistory(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var payments []models.Payment
	config.DB.Where("user_id = ?", user.ID).Find(&payments)

	c.JSON(http.StatusOK, payments)
}

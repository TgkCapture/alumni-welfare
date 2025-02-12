package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/TgkCapture/alumni-welfare/config"
	"github.com/TgkCapture/alumni-welfare/models"
	"github.com/gin-gonic/gin"
)

type PaymentRequest struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
	Month  int    `json:"month"`
}

type PayChanguResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

// MakePayment - Initiate a Mobile Money Payment
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

	// Call PayChangu API
	payChanguURL := os.Getenv("PAYCHANGU_BASE_URL") + "/mobile-money/payments/initialize"
	payChanguPayload := map[string]interface{}{
		"mobile_money_operator_ref_id": "20be6c20-adeb-4b5b-a7ba-0769820df4fb",
	}
	payloadBytes, _ := json.Marshal(payChanguPayload)

	req, err := http.NewRequest("POST", payChanguURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request creation failed"})
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("PAYCHANGU_SECRET_KEY"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment request failed"})
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var payChanguResponse PayChanguResponse
	json.Unmarshal(body, &payChanguResponse)

	// Save payment record
	payment := models.Payment{
		UserID:        user.ID,
		Name:          request.Name,
		Amount:        request.Amount,
		Month:         request.Month,
		TransactionID: payChanguResponse.TransactionID,
		Status:        "pending",
	}
	config.DB.Create(&payment)

	c.JSON(http.StatusOK, gin.H{"message": "Payment initiated", "transaction_id": payChanguResponse.TransactionID})
}

// Webhook handler for PayChangu payment updates
func PaymentWebhook(c *gin.Context) {
	var webhookEvent struct {
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"` // "successful" or "failed"
	}

	if err := c.ShouldBindJSON(&webhookEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Update payment record in database
	var payment models.Payment
	if err := config.DB.Where("transaction_id = ?", webhookEvent.TransactionID).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	payment.Status = webhookEvent.Status
	config.DB.Save(&payment)

	c.JSON(http.StatusOK, gin.H{"message": "Payment updated", "status": webhookEvent.Status})
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

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/TgkCapture/alumni-welfare/config"
	"github.com/TgkCapture/alumni-welfare/models"
	"github.com/TgkCapture/alumni-welfare/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type PaymentVerificationResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
	Details       string `json:"details"`
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

	// Retrieve user from the database
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Fetch Mobile Money Operator ID
	operatorID, err := services.GetOperatorRefID(user.MobileNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find mobile money operator"})
		return
	}

	chargeID := uuid.New().String()

	payChanguPayload := map[string]interface{}{
		"mobile_money_operator_ref_id": operatorID,
		"mobile":                       user.MobileNumber,
		"email":                        user.Email,
		"first_name":                   user.FirstName,
		"last_name":                    user.LastName,
		"amount":                       request.Amount,
		"charge_id":                    chargeID,
	}

	payloadBytes, err := json.Marshal(payChanguPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment payload"})
		return
	}

	// Send the request to PayChangu API
	payChanguURL := os.Getenv("PAYCHANGU_BASE_URL") + "/mobile-money/payments/initialize"
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

	// Parse the response from PayChangu
	body, _ := io.ReadAll(res.Body)
	var payChanguResponse PayChanguResponse
	json.Unmarshal(body, &payChanguResponse)

	// Save the payment record to the database
	payment := models.Payment{
		UserID:        user.ID,
		Name:          request.Name,
		Amount:        request.Amount,
		Month:         request.Month,
		TransactionID: payChanguResponse.TransactionID,
		Status:        "pending",
	}
	config.DB.Create(&payment)

	c.JSON(http.StatusOK, gin.H{"message": "Payment initiated", "transaction_id": payChanguResponse.TransactionID, "charge_id": chargeID})
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

// VerifyPayment - Verify Mobile Money
func VerifyPayment(c *gin.Context) {
	chargeID := c.Param("chargeId")

	if chargeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing chargeId"})
		return
	}

	verifyURL := fmt.Sprintf("https://api.paychangu.com/mobile-money/payments/%s/verify", chargeID)

	req, err := http.NewRequest("GET", verifyURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request creation failed"})
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("PAYCHANGU_SECRET_KEY"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification request failed"})
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var verificationResponse PaymentVerificationResponse
	if err := json.Unmarshal(body, &verificationResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse verification response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id": verificationResponse.TransactionID,
		"status":         verificationResponse.Status,
		"amount":         verificationResponse.Amount,
		"currency":       verificationResponse.Currency,
		"details":        verificationResponse.Details,
	})
}

// GetPaymentDetails - Retrieve detailed payment information
func GetPaymentDetails(c *gin.Context) {
	chargeID := c.Param("chargeId")

	if chargeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing chargeId"})
		return
	}

	detailsURL := fmt.Sprintf("https://api.paychangu.com/mobile-money/payments/%s/details", chargeID)

	req, err := http.NewRequest("GET", detailsURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request creation failed"})
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("PAYCHANGU_SECRET_KEY"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Details request failed"})
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var detailsResponse PaymentVerificationResponse
	if err := json.Unmarshal(body, &detailsResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse payment details response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_id": detailsResponse.TransactionID,
		"status":         detailsResponse.Status,
		"amount":         detailsResponse.Amount,
		"currency":       detailsResponse.Currency,
		"details":        detailsResponse.Details,
	})
}

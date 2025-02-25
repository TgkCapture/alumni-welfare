package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/TgkCapture/alumni-welfare/config"
	"github.com/TgkCapture/alumni-welfare/models"
	"github.com/TgkCapture/alumni-welfare/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Amount    int    `json:"amount"`
	Month     int    `json:"month"`
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
		FirstName:     request.FirstName,
		LastName:      request.LastName,
		Amount:        request.Amount,
		Month:         request.Month,
		TransactionID: payChanguResponse.TransactionID,
		Status:        "pending",
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: Failed to save payment"})
		return
	}

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

	var payment models.Payment
	err := config.DB.Where("transaction_id = ?", webhookEvent.TransactionID).First(&payment).Error

	if err != nil {
		// If no existing record, create a new one
		if errors.Is(err, gorm.ErrRecordNotFound) {
			payment = models.Payment{
				TransactionID: webhookEvent.TransactionID,
				Status:        webhookEvent.Status,
			}
			if err := config.DB.Create(&payment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment record"})
			return
		}
	} else {
		// Update existing record
		payment.Status = webhookEvent.Status
		if err := config.DB.Save(&payment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
			return
		}
	}

	// Generate transaction report
	reportPath, err := services.GenerateTransactionReport(payment)
	if err == nil {
		payment.ReportPath = reportPath
		config.DB.Save(&payment) // Save the report path
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Payment updated",
		"status":     webhookEvent.Status,
		"report_url": reportPath,
	})
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
	if err := config.DB.Where("user_id = ?", user.ID).Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment history"})
		return
	}

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

	req, _ := http.NewRequest("GET", verifyURL, nil)
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
	json.Unmarshal(body, &verificationResponse)

	// Update payment status in database if successful
	if verificationResponse.Status == "successful" {
		config.DB.Model(&models.Payment{}).
			Where("transaction_id = ?", verificationResponse.TransactionID).
			Update("status", "successful")
	}

	c.JSON(http.StatusOK, verificationResponse)
}

// GetPaymentDetails - Retrieve detailed payment information
func GetPaymentDetails(c *gin.Context) {
	chargeID := c.Param("chargeId")

	if chargeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing chargeId"})
		return
	}

	detailsURL := fmt.Sprintf("https://api.paychangu.com/mobile-money/payments/%s/details", chargeID)
	req, _ := http.NewRequest("GET", detailsURL, nil)
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
	json.Unmarshal(body, &detailsResponse)

	c.JSON(http.StatusOK, detailsResponse)
}

// serves the transaction report file
func GetTransactionReport(c *gin.Context) {
	transactionID := c.Param("transactionId")

	var payment models.Payment
	if err := config.DB.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if payment.ReportPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not available"})
		return
	}

	c.File(payment.ReportPath)
}

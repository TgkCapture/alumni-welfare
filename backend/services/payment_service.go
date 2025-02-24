package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/TgkCapture/alumni-welfare/models"
	"github.com/TgkCapture/alumni-welfare/utils"
	"github.com/gin-gonic/gin"
)

func GetMobileOperator(c *gin.Context) {
	url := os.Getenv("PAYCHANGU_BASE_URL") + "/mobile-money"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Add("accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mobile operators"})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Parse JSON response
	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, responseData)
}

func GetOperatorRefID(mobile string) (string, error) {
	url := os.Getenv("PAYCHANGU_BASE_URL") + "/mobile-money"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var responseData struct {
		Data []utils.Operator `json:"data"`
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", err
	}

	// Determine the short code based on the mobile number prefix
	var shortCode string
	switch {
	case strings.HasPrefix(mobile, "088"):
		shortCode = "tnm"
	case strings.HasPrefix(mobile, "099"), strings.HasPrefix(mobile, "098"):
		shortCode = "airtel"
	default:
		return "", nil
	}

	return utils.FindOperatorRefID(responseData.Data, shortCode), nil
}

// GenerateTransactionReport creates a CSV report of the transaction
func GenerateTransactionReport(payment models.Payment) (string, error) {
	reportDir := "reports"
	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		os.Mkdir(reportDir, os.ModePerm)
	}

	// Report filename
	filename := fmt.Sprintf("%s/transaction_%s.csv", reportDir, payment.TransactionID)
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write CSV headers
	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Transaction ID", "User ID", "Name", "Amount", "Month", "Status", "Timestamp"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	// Write transaction data
	data := []string{
		payment.TransactionID,
		fmt.Sprintf("%d", payment.UserID),
		payment.Name,
		fmt.Sprintf("%d", payment.Amount),
		fmt.Sprintf("%d", payment.Month),
		payment.Status,
		time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := writer.Write(data); err != nil {
		return "", err
	}

	return filename, nil
}

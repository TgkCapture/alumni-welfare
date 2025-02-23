package services

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

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

package services

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

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
		Data []struct {
			RefID     string `json:"ref_id"`
			ShortCode string `json:"short_code"`
			Name      string `json:"name"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", err
	}

	// Check mobile number prefix
	if strings.HasPrefix(mobile, "088") {
		return FindOperatorRefID(responseData.Data, "tnm"), nil
	} else if strings.HasPrefix(mobile, "099") || strings.HasPrefix(mobile, "098") {
		return FindOperatorRefID(responseData.Data, "airtel"), nil
	}

	return "", nil
}

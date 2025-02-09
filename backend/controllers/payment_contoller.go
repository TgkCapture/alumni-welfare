package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentRequest struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

func MakePayment(c *gin.Context) {
	var request PaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// TODO: integrate a real payment gateway
	fmt.Printf("Payment received: %s - %d\n", request.Name, request.Amount)

	c.JSON(http.StatusOK, gin.H{"message": "Payment received successfully!"})
}

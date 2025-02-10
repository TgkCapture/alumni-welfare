package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	UserID        uint   `json:"user_id"`
	Name          string `json:"name"`
	Amount        int    `json:"amount"`
	Month         int    `json:"month"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

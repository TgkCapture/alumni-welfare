package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	UserID        uint   `json:"user_id" gorm:"index"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Amount        int    `json:"amount"`
	Month         int    `json:"month"`
	TransactionID string `json:"transaction_id" gorm:"unique"`
	Status        string `json:"status" gorm:"index"`
	ReportPath    string `json:"report_path"`
}

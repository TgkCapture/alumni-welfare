package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	UserID        uint   `json:"user_id" gorm:"index"`
	Name          string `json:"name"`
	Amount        int    `json:"amount"`
	Month         int    `json:"month"`
	TransactionID string `json:"transaction_id" gorm:"unique"`
	Status        string `json:"status" gorm:"index"`
	ReportPath    string `json:"report_path"`
}

package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Month  string `json:"month"`
	Amount int    `json:"amount"`
}

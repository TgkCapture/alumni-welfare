package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Name         string    `json:"name"`
	Email        string    `json:"email" gorm:"unique"`
	MobileNumber string    `json:"mobile_number"`
	Password     string    `json:"-"`
	Payments     []Payment `gorm:"foreignKey:UserID"`
}

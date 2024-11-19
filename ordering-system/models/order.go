package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderNumber string    `json:"orderNumber"`
	CustomerID  uint      `json:"customerId"`
	Customer    Customer  `json:"customer"`
	Products    []Product `gorm:"many2many:order_products;" json:"products"`
}


package database

import (
	"ordering-system/models"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&models.Customer{}, &models.Order{}, &models.Product{})
}

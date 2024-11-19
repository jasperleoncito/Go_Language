package database

import (
	"log"
	"ordering-system/models"

	"gorm.io/gorm"
)

func SeedDB(db *gorm.DB) {
	// Seed Customers
	customers := []models.Customer{
		{Name: "John Doe", Email: "john@example.com"},
		{Name: "Jane Doe", Email: "jane@example.com"},
	}
	for _, customer := range customers {
		db.FirstOrCreate(&customer)
	}

	// Seed Products
	products := []models.Product{
		{Name: "Laptop", Price: 1000},
		{Name: "Mouse", Price: 25},
	}
	for _, product := range products {
		db.FirstOrCreate(&product)
	}

	// Fetch the products to associate them with orders
	var laptop, mouse models.Product
	if err := db.First(&laptop, "name = ?", "Laptop").Error; err != nil {
		log.Fatalf("Failed to find product Laptop: %v", err)
	}
	if err := db.First(&mouse, "name = ?", "Mouse").Error; err != nil {
		log.Fatalf("Failed to find product Mouse: %v", err)
	}

	// Seed Orders with products
	orders := []models.Order{
		{
			OrderNumber: "ORD123",
			CustomerID:  1,
			Products:    []models.Product{laptop, mouse},
		},
		{
			OrderNumber: "ORD124",
			CustomerID:  2,
		},
	}
	for _, order := range orders {
		if err := db.Create(&order).Error; err != nil {
			log.Fatalf("Failed to create order: %v", err)
		}
	}
}

package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ordering-system/models"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Utility functions for pagination, filtering, and sorting
func paginate(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}
	return page, pageSize
}

func sort(r *http.Request) (string, string) {
	sortQuery := r.URL.Query().Get("sort")
	sortParts := strings.Split(sortQuery, ":")
	column, direction := "id", "asc" // Default sorting by ID
	if len(sortParts) == 2 {
		column, direction = sortParts[0], sortParts[1]
	}
	return column, direction
}

func filter(r *http.Request, db *gorm.DB) *gorm.DB {
	filterName := r.URL.Query().Get("name")
	if filterName != "" {
		db = db.Where("name ILIKE ?", "%"+filterName+"%")
	}
	return db
}

// CreateCustomer creates a new customer
func CreateCustomer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customer models.Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.Create(&customer).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(customer)
	}
}

// UpdateCustomer updates an existing customer
func UpdateCustomer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var customer models.Customer
		if err := db.First(&customer, id).Error; err != nil {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.Save(&customer).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(customer)
	}
}

// GetCustomer retrieves a specific customer by ID
func GetCustomer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var customer models.Customer
		if err := db.Preload("Orders").First(&customer, id).Error; err != nil {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(customer)
	}
}

// CreateOrder creates a new order
func CreateOrder(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order models.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.Create(&order).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}

// UpdateOrder updates an existing order
func UpdateOrder(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var order models.Order
		if err := db.First(&order, id).Error; err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.Save(&order).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(order)
	}
}

func GetOrder(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var order models.Order
		if err := db.Preload("Products").Preload("Products.OrderProducts").Preload("Customer").First(&order, id).Error; err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(order)
	}
}

// CreateProduct creates a new product
func CreateProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.Create(&product).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}

// UpdateProduct updates an existing product
func UpdateProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var product models.Product
		if err := db.First(&product, id).Error; err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := db.Save(&product).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(product)
	}
}

// GetProduct retrieves a specific product by ID
func GetProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var product models.Product
		if err := db.First(&product, id).Error; err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(product)
	}
}

// ListCustomers lists all customers with pagination, filtering, and sorting
func ListCustomers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, pageSize := paginate(r)
		column, direction := sort(r)

		var customers []models.Customer
		db = filter(r, db)
		db.Order(clause.OrderByColumn{
			Column: clause.Column{Name: column},
			Desc:   direction == "desc",
		}).Offset((page - 1) * pageSize).Limit(pageSize).Find(&customers)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customers)
	}
}

// ListOrders lists all orders with pagination and sorting
func ListOrders(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, pageSize := paginate(r)
		column, direction := sort(r)

		var orders []models.Order
		db.Order(clause.OrderByColumn{
			Column: clause.Column{Name: column},
			Desc:   direction == "desc",
		}).Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
}

// ListProducts lists all products with pagination, filtering, and sorting
func ListProducts(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, pageSize := paginate(r)
		column, direction := sort(r)

		var products []models.Product
		db = filter(r, db)
		db.Order(clause.OrderByColumn{
			Column: clause.Column{Name: column},
			Desc:   direction == "desc",
		}).Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

// Routes maps all the endpoints
func Routes(db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	// Customer Routes
	r.Get("/customers", ListCustomers(db))
	r.Post("/customers", CreateCustomer(db))
	r.Get("/customers/{id}", GetCustomer(db))
	r.Put("/customers/{id}", UpdateCustomer(db))

	// Order Routes
	r.Get("/orders", ListOrders(db))
	r.Post("/orders", CreateOrder(db))
	r.Get("/orders/{id}", GetOrder(db))
	r.Put("/orders/{id}", UpdateOrder(db))

	// Product Routes
	r.Get("/products", ListProducts(db))
	r.Post("/products", CreateProduct(db))
	r.Get("/products/{id}", GetProduct(db))
	r.Put("/products/{id}", UpdateProduct(db))

	return r
}

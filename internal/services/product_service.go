package services

import (
	"errors"
	"time"

	"github.com/Bitovi/example-go-server/internal/models"
)

var (
	// ErrProductNotFound is returned when a product is not found
	ErrProductNotFound = errors.New("product not found")
	
	// Mock product data
	mockProducts = []models.Product{
		{
			ID:          "550e8400-e29b-41d4-a716-446655440000",
			Name:        "Laptop",
			Description: "High-performance laptop for professionals",
			Price:       1299.99,
			Category:    "Electronics",
			InStock:     true,
			CreatedAt:   time.Now().AddDate(0, -3, 0),
			UpdatedAt:   time.Now().AddDate(0, -1, 0),
		},
		{
			ID:          "550e8400-e29b-41d4-a716-446655440001",
			Name:        "Wireless Mouse",
			Description: "Ergonomic wireless mouse with precision tracking",
			Price:       29.99,
			Category:    "Electronics",
			InStock:     true,
			CreatedAt:   time.Now().AddDate(0, -2, 0),
			UpdatedAt:   time.Now().AddDate(0, 0, -5),
		},
		{
			ID:          "550e8400-e29b-41d4-a716-446655440002",
			Name:        "Desk Lamp",
			Description: "LED desk lamp with adjustable brightness",
			Price:       49.99,
			Category:    "Office",
			InStock:     false,
			CreatedAt:   time.Now().AddDate(0, -1, 0),
			UpdatedAt:   time.Now().AddDate(0, 0, -2),
		},
		{
			ID:          "550e8400-e29b-41d4-a716-446655440003",
			Name:        "Notebook",
			Description: "Premium leather-bound notebook",
			Price:       19.99,
			Category:    "Office",
			InStock:     true,
			CreatedAt:   time.Now().AddDate(0, -4, 0),
			UpdatedAt:   time.Now().AddDate(0, -1, -10),
		},
		{
			ID:          "550e8400-e29b-41d4-a716-446655440004",
			Name:        "Coffee Maker",
			Description: "Programmable coffee maker with timer",
			Price:       79.99,
			Category:    "Kitchen",
			InStock:     true,
			CreatedAt:   time.Now().AddDate(0, -5, 0),
			UpdatedAt:   time.Now().AddDate(0, -2, 0),
		},
	}
)

// ProductService handles business logic for products
type ProductService struct{}

// NewProductService creates a new product service
func NewProductService() *ProductService {
	return &ProductService{}
}

// ListProducts returns a list of products with optional limit
func (s *ProductService) ListProducts(limit int) ([]models.Product, int) {
	total := len(mockProducts)
	
	if limit <= 0 || limit > len(mockProducts) {
		limit = len(mockProducts)
	}
	
	products := make([]models.Product, limit)
	copy(products, mockProducts[:limit])
	
	return products, total
}

// GetProductByID returns a product by its ID
func (s *ProductService) GetProductByID(id string) (*models.Product, error) {
	for _, product := range mockProducts {
		if product.ID == id {
			// Return a copy to prevent modification
			p := product
			return &p, nil
		}
	}
	
	return nil, ErrProductNotFound
}

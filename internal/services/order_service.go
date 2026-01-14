package services

import (
	"errors"
	"time"

	"github.com/Bitovi/example-go-server/internal/models"
	"github.com/google/uuid"
)

var (
	// ErrOrderNotFound is returned when an order is not found
	ErrOrderNotFound = errors.New("order not found")

	// orderUserMap tracks which user owns which order
	orderUserMap = map[string]string{
		"650e8400-e29b-41d4-a716-446655440000": "750e8400-e29b-41d4-a716-446655440000", // johndoe
		"650e8400-e29b-41d4-a716-446655440001": "750e8400-e29b-41d4-a716-446655440000", // johndoe
		"650e8400-e29b-41d4-a716-446655440002": "750e8400-e29b-41d4-a716-446655440001", // janedoe
	}

	// Mock order data
	mockOrders = []models.Order{
		{
			ID: "650e8400-e29b-41d4-a716-446655440000",
			Products: []models.OrderProduct{
				{
					ProductID: "550e8400-e29b-41d4-a716-446655440000", // Laptop
					Quantity:  1,
				},
				{
					ProductID: "550e8400-e29b-41d4-a716-446655440001", // Wireless Mouse
					Quantity:  2,
				},
			},
			TotalPrice:           1359.97,
			AccruedLoyaltyPoints: 135, // 1359.97 / 10 = 135 points
			OrderDate:            time.Now().AddDate(0, 0, -5),
			Status:               models.OrderStatusPending,
		},
		{
			ID: "650e8400-e29b-41d4-a716-446655440001",
			Products: []models.OrderProduct{
				{
					ProductID: "550e8400-e29b-41d4-a716-446655440002", // Desk Lamp
					Quantity:  3,
				},
			},
			TotalPrice:           149.97,
			AccruedLoyaltyPoints: 14, // 149.97 / 10 = 14 points
			OrderDate:            time.Now().AddDate(0, 0, -3),
			Status:               models.OrderStatusShipped,
		},
		{
			ID: "650e8400-e29b-41d4-a716-446655440002",
			Products: []models.OrderProduct{
				{
					ProductID: "550e8400-e29b-41d4-a716-446655440003", // Notebook
					Quantity:  5,
				},
				{
					ProductID: "550e8400-e29b-41d4-a716-446655440004", // Coffee Maker
					Quantity:  1,
				},
			},
			TotalPrice:           179.94,
			AccruedLoyaltyPoints: 17, // 179.94 / 10 = 17 points
			OrderDate:            time.Now().AddDate(0, 0, -1),
			Status:               models.OrderStatusProcessing,
		},
	}
)

// ResetOrderMockData resets the mock order data to its initial state
// This should be called in test setup to ensure test isolation
func ResetOrderMockData() {
	orderUserMap = map[string]string{
		"650e8400-e29b-41d4-a716-446655440000": "750e8400-e29b-41d4-a716-446655440000", // johndoe
		"650e8400-e29b-41d4-a716-446655440001": "750e8400-e29b-41d4-a716-446655440000", // johndoe
		"650e8400-e29b-41d4-a716-446655440002": "750e8400-e29b-41d4-a716-446655440001", // janedoe
	}

	mockOrders = []models.Order{
		{
			ID: "650e8400-e29b-41d4-a716-446655440000",
			Products: []models.OrderProduct{
				{
					ProductID: "550e8400-e29b-41d4-a716-446655440000", // Laptop
					Quantity:  1,
				},
				{
					ProductID: "550e8400-e29b-41d4-a716-446655440001", // Wireless Mouse
					Quantity:  2,
				},
			},
			TotalPrice:           1359.97,
			AccruedLoyaltyPoints: 135, // 1359.97 / 10 = 135 points
			OrderDate:            time.Now().AddDate(0, 0, -5),
			Status:               models.OrderStatusPending,
		},
		{
			ID: "650e8400-e29b-41d4-a716-446655440001",
			Products: []models.OrderProduct{
				{
					ProductID: "550e8400-e29b-41d4-a716-446655440002", // Desk Lamp
					Quantity:  3,
				},
			},
			TotalPrice:           149.97,
			AccruedLoyaltyPoints: 14, // 149.97 / 10 = 14 points
			OrderDate:            time.Now().AddDate(0, 0, -3),
			Status:               models.OrderStatusShipped,
		},
		{
			ID: "650e8400-e29b-41d4-a716-446655440002",
			Products: []models.OrderProduct{
				{
					ProductID: "550e8400-e29b-41d4-a716-446655440003", // Notebook
					Quantity:  5,
				},
				{
					ProductID: "550e8400-e29b-41d4-a716-446655440004", // Coffee Maker
					Quantity:  1,
				},
			},
			TotalPrice:           179.94,
			AccruedLoyaltyPoints: 17, // 179.94 / 10 = 17 points
			OrderDate:            time.Now().AddDate(0, 0, -1),
			Status:               models.OrderStatusProcessing,
		},
	}
}

// GetMockOrders returns a copy of mock orders for cross-service access
func GetMockOrders() []models.Order {
	orders := make([]models.Order, len(mockOrders))
	copy(orders, mockOrders)
	return orders
}

// GetMockOrdersReference returns a reference to the actual mock orders slice
// This allows other services to modify orders directly (e.g., cancelling on user deletion)
func GetMockOrdersReference() []models.Order {
	return mockOrders
}

// UpdateMockOrderStatus updates the status of an order at the given index
// This is used by UserService to cancel pending orders when deleting a user
func UpdateMockOrderStatus(index int, status models.OrderStatus) {
	if index >= 0 && index < len(mockOrders) {
		mockOrders[index].Status = status
	}
}

// OrderService handles business logic for orders
type OrderService struct {}

// NewOrderService creates a new order service
func NewOrderService() *OrderService {
	return &OrderService{}
}

// ListOrders returns a list of all orders
func (s *OrderService) ListOrders() ([]models.Order, int) {
	total := len(mockOrders)

	// Return a copy to prevent modification
	orders := make([]models.Order, len(mockOrders))
	copy(orders, mockOrders)

	return orders, total
}

// GetOrderByID returns an order by its ID
func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
	for _, order := range mockOrders {
		if order.ID == id {
			// Return a copy to prevent modification
			o := order
			return &o, nil
		}
	}

	return nil, ErrOrderNotFound
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(userID string, products []models.OrderProduct) (*models.Order, error) {
	if len(products) == 0 {
		return nil, errors.New("order must contain at least one product")
	}

	// Calculate total price
	// In a real implementation, this would look up product prices from the product service
	totalPrice := 0.0
	for _, orderProduct := range products {
		// For now, use a placeholder price calculation
		// TODO: Integrate with PRODUCT_SERVICE_URL to fetch actual product prices
		totalPrice += 10.0 * float64(orderProduct.Quantity) // Placeholder: $10 per item
	}

	// Calculate accrued loyalty points: 1 point per $10 spent (rounded down)
	accruedPoints := int(totalPrice / 10.0)

	// Generate new order with proper UUID
	orderID := uuid.New().String()
	newOrder := models.Order{
		ID:                   orderID,
		Products:             products,
		TotalPrice:           totalPrice,
		AccruedLoyaltyPoints: accruedPoints,
		OrderDate:            time.Now(),
		Status:               models.OrderStatusPending,
	}

	// If userId is provided, track the order-user relationship
	if userID != "" {
		orderUserMap[orderID] = userID
	}

	// Add to mock orders
	mockOrders = append(mockOrders, newOrder)

	return &newOrder, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(orderID string, status models.OrderStatus) (*models.Order, error) {
	for i, order := range mockOrders {
		if order.ID == orderID {
			mockOrders[i].Status = status
			return &mockOrders[i], nil
		}
	}

	return nil, ErrOrderNotFound
}

// UpdateOrderProducts updates the products in an order (only for PENDING orders)
// For each product in the input:
// - If quantity > 0: adds the quantity to existing product (or creates new product)
// - If quantity < 0: subtracts the quantity from existing product (removes if result <= 0)
// - If quantity = 0: does nothing
func (s *OrderService) UpdateOrderProducts(orderID string, products []models.OrderProduct) (*models.Order, error) {
	for i, order := range mockOrders {
		if order.ID == orderID {
			// Only allow updating products for pending orders
			if order.Status != models.OrderStatusPending {
				return nil, errors.New("can only update products for pending orders")
			}

			// Create a map of existing products for quick lookup
			existingProducts := make(map[string]models.OrderProduct)
			for _, product := range order.Products {
				existingProducts[product.ProductID] = product
			}

			// Process each product in the request
			for _, product := range products {
				if product.Quantity == 0 {
					// Do nothing
					continue
				}

				existing, exists := existingProducts[product.ProductID]
				if exists {
					// Product already exists - add or subtract quantity
					newQuantity := existing.Quantity + product.Quantity
					if newQuantity <= 0 {
						// Remove the product if quantity becomes 0 or negative
						delete(existingProducts, product.ProductID)
					} else {
						// Update the quantity
						existing.Quantity = newQuantity
						existingProducts[product.ProductID] = existing
					}
				} else if product.Quantity > 0 {
					// New product with positive quantity - add it
					existingProducts[product.ProductID] = product
				}
				// If product doesn't exist and quantity is negative, ignore it
			}

			// Convert map back to slice
			updatedProducts := make([]models.OrderProduct, 0, len(existingProducts))
			for _, product := range existingProducts {
				updatedProducts = append(updatedProducts, product)
			}

			// Recalculate total price
			// In a real implementation, this would look up product prices from the product service
			totalPrice := 0.0
			for _, orderProduct := range updatedProducts {
				// For now, use a placeholder price calculation
				// TODO: Integrate with PRODUCT_SERVICE_URL to fetch actual product prices
				totalPrice += 10.0 * float64(orderProduct.Quantity) // Placeholder: $10 per item
			}

			// Recalculate accrued loyalty points: 1 point per $10 spent (rounded down)
			accruedPoints := int(totalPrice / 10.0)

			// Update the order
			mockOrders[i].Products = updatedProducts
			mockOrders[i].TotalPrice = totalPrice
			mockOrders[i].AccruedLoyaltyPoints = accruedPoints

			return &mockOrders[i], nil
		}
	}

	return nil, ErrOrderNotFound
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(orderID string) (*models.Order, error) {
	return s.UpdateOrderStatus(orderID, models.OrderStatusCanceled)
}

// SubmitOrder submits a pending order for processing
func (s *OrderService) SubmitOrder(orderID string) (*models.Order, error) {
	for i, order := range mockOrders {
		if order.ID == orderID {
			if order.Status != models.OrderStatusPending {
				return nil, errors.New("only pending orders can be submitted")
			}
			mockOrders[i].Status = models.OrderStatusProcessing

			// Award loyalty points to the user (1 point per $10 spent)
			// TODO: Integrate with LOYALTY_SERVICE_URL to award loyalty points
			if userID, ok := orderUserMap[orderID]; ok {
				_ = userID // Placeholder: would call loyalty service here
				_ = order.AccruedLoyaltyPoints
			}

			return &mockOrders[i], nil
		}
	}

	return nil, ErrOrderNotFound
}

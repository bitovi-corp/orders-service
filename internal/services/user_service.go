package services

import (
	"errors"
	"time"

	"github.com/Bitovi/example-go-server/internal/models"
	"github.com/google/uuid"
)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")

	// Mock user data
	mockUsers = []models.User{
		{
			ID:            "750e8400-e29b-41d4-a716-446655440000",
			Username:      "johndoe",
			Email:         "john.doe@example.com",
			Firstname:     "John",
			Lastname:      "Doe",
			LoyaltyPoints: 1500,
			CreatedAt:     time.Now().AddDate(0, -6, 0),
			UpdatedAt:     time.Now().AddDate(0, -1, 0),
		},
		{
			ID:            "750e8400-e29b-41d4-a716-446655440001",
			Username:      "janedoe",
			Email:         "jane.doe@example.com",
			Firstname:     "Jane",
			Lastname:      "Doe",
			LoyaltyPoints: 2300,
			CreatedAt:     time.Now().AddDate(0, -4, 0),
			UpdatedAt:     time.Now().AddDate(0, 0, -10),
		},
		{
			ID:            "750e8400-e29b-41d4-a716-446655440002",
			Username:      "bobsmith",
			Email:         "bob.smith@example.com",
			Firstname:     "Bob",
			Lastname:      "Smith",
			LoyaltyPoints: 500,
			CreatedAt:     time.Now().AddDate(0, -8, 0),
			UpdatedAt:     time.Now().AddDate(0, -2, 0),
		},
	}

	// Map users to their orders
	userOrdersMap = map[string][]string{
		"750e8400-e29b-41d4-a716-446655440000": {"650e8400-e29b-41d4-a716-446655440000", "650e8400-e29b-41d4-a716-446655440001"},
		"750e8400-e29b-41d4-a716-446655440001": {"650e8400-e29b-41d4-a716-446655440002"},
		"750e8400-e29b-41d4-a716-446655440002": {},
	}
)

// ResetUserMockData resets the mock user data to its initial state
// This should be called in test setup to ensure test isolation
func ResetUserMockData() {
	mockUsers = []models.User{
		{
			ID:            "750e8400-e29b-41d4-a716-446655440000",
			Username:      "johndoe",
			Email:         "john.doe@example.com",
			Firstname:     "John",
			Lastname:      "Doe",
			LoyaltyPoints: 1500,
			CreatedAt:     time.Now().AddDate(0, -6, 0),
			UpdatedAt:     time.Now().AddDate(0, -1, 0),
		},
		{
			ID:            "750e8400-e29b-41d4-a716-446655440001",
			Username:      "janedoe",
			Email:         "jane.doe@example.com",
			Firstname:     "Jane",
			Lastname:      "Doe",
			LoyaltyPoints: 2300,
			CreatedAt:     time.Now().AddDate(0, -4, 0),
			UpdatedAt:     time.Now().AddDate(0, 0, -10),
		},
		{
			ID:            "750e8400-e29b-41d4-a716-446655440002",
			Username:      "bobsmith",
			Email:         "bob.smith@example.com",
			Firstname:     "Bob",
			Lastname:      "Smith",
			LoyaltyPoints: 500,
			CreatedAt:     time.Now().AddDate(0, -8, 0),
			UpdatedAt:     time.Now().AddDate(0, -2, 0),
		},
	}

	userOrdersMap = map[string][]string{
		"750e8400-e29b-41d4-a716-446655440000": {"650e8400-e29b-41d4-a716-446655440000", "650e8400-e29b-41d4-a716-446655440001"},
		"750e8400-e29b-41d4-a716-446655440001": {"650e8400-e29b-41d4-a716-446655440002"},
		"750e8400-e29b-41d4-a716-446655440002": {},
	}
}

// UserService handles business logic for users
type UserService struct{}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{}
}

// AddOrderToUser adds an order ID to a user's order list
// This is called from OrderService when creating an order
func AddOrderToUser(userID, orderID string) {
	if orders, exists := userOrdersMap[userID]; exists {
		userOrdersMap[userID] = append(orders, orderID)
	} else {
		userOrdersMap[userID] = []string{orderID}
	}
}

// GetUserByID returns a user by their ID
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	for _, user := range mockUsers {
		if user.ID == id {
			// Return a copy to prevent modification
			u := user
			return &u, nil
		}
	}

	return nil, ErrUserNotFound
}

// GetUserWithOrders returns a user with their associated orders
func (s *UserService) GetUserWithOrders(userID string) (*models.UserOrders, error) {
	// Get user
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Get user's order IDs
	orderIDs, exists := userOrdersMap[userID]
	if !exists {
		// User exists but has no orders
		return &models.UserOrders{
			User:   *user,
			Orders: []models.Order{},
		}, nil
	}

	// Get all orders from order service
	allOrders := GetMockOrders()

	// Filter orders for this user
	userOrders := make([]models.Order, 0)
	for _, order := range allOrders {
		for _, orderID := range orderIDs {
			if order.ID == orderID {
				userOrders = append(userOrders, order)
				break
			}
		}
	}

	return &models.UserOrders{
		User:   *user,
		Orders: userOrders,
	}, nil
}

// GetUserLoyaltyPoints returns the loyalty points for a user
func (s *UserService) GetUserLoyaltyPoints(userID string) (int, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return 0, err
	}
	
	return user.LoyaltyPoints, nil
}

// RedeemUserLoyaltyPoints redeems loyalty points for a user
func (s *UserService) RedeemUserLoyaltyPoints(userID string, pointsToRedeem int) (int, error) {
	if pointsToRedeem < 1 {
		return 0, errors.New("points to redeem must be at least 1")
	}

	for i, user := range mockUsers {
		if user.ID == userID {
			if user.LoyaltyPoints < pointsToRedeem {
				return 0, errors.New("insufficient loyalty points")
			}
			
			mockUsers[i].LoyaltyPoints -= pointsToRedeem
			return mockUsers[i].LoyaltyPoints, nil
		}
	}

	return 0, ErrUserNotFound
}

// AwardLoyaltyPoints awards loyalty points to a user
func (s *UserService) AwardLoyaltyPoints(userID string, points int) error {
	if points < 0 {
		return errors.New("points to award must be non-negative")
	}

	for i, user := range mockUsers {
		if user.ID == userID {
			mockUsers[i].LoyaltyPoints += points
			return nil
		}
	}

	return ErrUserNotFound
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, email, firstname, lastname string) (*models.User, error) {
	// Validate username length
	if len(username) < 3 || len(username) > 30 {
		return nil, errors.New("username must be between 3 and 30 characters")
	}

	// Basic email validation (format check is handled by handler)
	if email == "" {
		return nil, errors.New("email is required")
	}

	// Check if username already exists
	for _, user := range mockUsers {
		if user.Username == username {
			return nil, errors.New("username already exists")
		}
		if user.Email == email {
			return nil, errors.New("email already exists")
		}
	}

	// Create new user
	newUser := models.User{
		ID:            uuid.New().String(),
		Username:      username,
		Email:         email,
		Firstname:     firstname,
		Lastname:      lastname,
		LoyaltyPoints: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Add to mock users
	mockUsers = append(mockUsers, newUser)

	// Initialize empty orders array for new user
	userOrdersMap[newUser.ID] = []string{}

	return &newUser, nil
}

// DeleteUser deletes a user by their ID
// Automatically cancels all PENDING orders for the user before deletion
func (s *UserService) DeleteUser(userID string) error {
	// Find the user index
	userIndex := -1
	for i, user := range mockUsers {
		if user.ID == userID {
			userIndex = i
			break
		}
	}

	if userIndex == -1 {
		return ErrUserNotFound
	}

	// Cancel all PENDING orders for this user
	orderIDs, exists := userOrdersMap[userID]
	if exists {
		// Access mockOrders directly from order_service to update orders
		allOrders := GetMockOrdersReference()
		for _, orderID := range orderIDs {
			for i := range allOrders {
				if allOrders[i].ID == orderID && allOrders[i].Status == models.OrderStatusPending {
					// Cancel the order by updating its status directly
					allOrders[i].Status = models.OrderStatusCanceled
				}
			}
		}
	}

	// Remove user from slice
	mockUsers = append(mockUsers[:userIndex], mockUsers[userIndex+1:]...)

	// Remove user from userOrdersMap
	delete(userOrdersMap, userID)

	return nil
}

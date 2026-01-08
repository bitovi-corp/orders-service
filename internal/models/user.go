package models

import (
	"time"
)

// User represents a user as defined in api/openapi.yaml
type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Firstname     string    `json:"firstname,omitempty"`
	Lastname      string    `json:"lastname,omitempty"`
	LoyaltyPoints int       `json:"loyaltyPoints,omitempty"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
}

// UserOrders represents a user with their associated orders
type UserOrders struct {
	User   User    `json:"user"`
	Orders []Order `json:"orders"`
}

package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Bitovi/example-go-server/internal/services"
	"github.com/google/uuid"
)

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// GetUserWithOrders implements GET /user/{userId} endpoint as defined in api/openapi.yaml
func GetUserWithOrders(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Extract user ID from URL path
	// URL format: /user/{userId}
	path := strings.TrimPrefix(r.URL.Path, "/user/")
	userID := strings.Split(path, "/")[0]

	if userID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID", "User ID is required", "")
		return
	}

	// UUID format validation using google/uuid
	if _, err := uuid.Parse(userID); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format", "User ID must be a valid UUID")
		return
	}

	// Get user with orders from service
	userOrders, err := userService.GetUserWithOrders(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "The requested user could not be found", "")
			return
		}
		log.Printf("Error retrieving user with orders: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", "")
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userOrders); err != nil {
		log.Printf("Error encoding user orders response: %v", err)
	}
}

// GetUserLoyaltyPoints implements GET /user/{userId}/points endpoint as defined in api/openapi.yaml
func GetUserLoyaltyPoints(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/user/")
	pathParts := strings.Split(path, "/")
	userID := pathParts[0]

	if userID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID", "User ID is required", "")
		return
	}

	// UUID format validation using google/uuid
	if _, err := uuid.Parse(userID); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format", "User ID must be a valid UUID")
		return
	}

	// Get loyalty points from service
	points, err := userService.GetUserLoyaltyPoints(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "The requested user could not be found", "")
			return
		}
		log.Printf("Error retrieving user loyalty points: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", "")
		return
	}

	// Prepare response
	response := map[string]int{
		"loyaltyPoints": points,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding loyalty points response: %v", err)
	}
}

// RedeemUserLoyaltyPoints implements POST /user/{userId}/points endpoint as defined in api/openapi.yaml
func RedeemUserLoyaltyPoints(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/user/")
	pathParts := strings.Split(path, "/")
	userID := pathParts[0]

	if userID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID", "User ID is required", "")
		return
	}

	// UUID format validation using google/uuid
	if _, err := uuid.Parse(userID); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format", "User ID must be a valid UUID")
		return
	}

	// Parse request body
	var requestBody struct {
		Points int `json:"points"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body", err.Error())
		return
	}

	// Validate points
	if requestBody.Points < 1 {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_POINTS", "Points to redeem must be at least 1", "")
		return
	}

	// Redeem points from service
	remainingPoints, err := userService.RedeemUserLoyaltyPoints(userID, requestBody.Points)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "The requested user could not be found", "")
			return
		}
		writeErrorResponse(w, http.StatusBadRequest, "REDEMPTION_FAILED", err.Error(), "")
		return
	}

	// Prepare response
	response := map[string]int{
		"remainingPoints": remainingPoints,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding redemption response: %v", err)
	}
}

// CreateUser implements POST /user endpoint as defined in api/openapi.yaml
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Parse request body
	var requestBody struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body", err.Error())
		return
	}

	// Validate required fields
	if requestBody.Username == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_INPUT", "Username is required", "")
		return
	}
	if requestBody.Email == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_INPUT", "Email is required", "")
		return
	}
	if requestBody.Firstname == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_INPUT", "Firstname is required", "")
		return
	}
	if requestBody.Lastname == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_INPUT", "Lastname is required", "")
		return
	}

	// Validate username length (3-30 characters)
	if len(requestBody.Username) < 3 || len(requestBody.Username) > 30 {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_INPUT", "Username must be between 3 and 30 characters", "")
		return
	}

	// Validate email format
	if !emailRegex.MatchString(requestBody.Email) {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid email format", "")
		return
	}

	// Create user via service
	user, err := userService.CreateUser(requestBody.Username, requestBody.Email, requestBody.Firstname, requestBody.Lastname)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "USER_CREATION_FAILED", err.Error(), "")
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding user creation response: %v", err)
	}
}

// DeleteUser implements DELETE /user/{userId} endpoint as defined in api/openapi.yaml
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Only allow DELETE method
	if r.Method != http.MethodDelete {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/user/")
	userID := strings.Split(path, "/")[0]

	if userID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID", "User ID is required", "")
		return
	}

	// UUID format validation using google/uuid
	if _, err := uuid.Parse(userID); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format", "User ID must be a valid UUID")
		return
	}

	// Delete user from service
	err := userService.DeleteUser(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "The requested user could not be found", "")
			return
		}
		log.Printf("Error deleting user: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", "")
		return
	}

	// Send 204 No Content response
	w.WriteHeader(http.StatusNoContent)
}

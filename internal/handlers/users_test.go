package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bitovi/example-go-server/internal/models"
	"github.com/Bitovi/example-go-server/internal/services"
)

// resetMockData should be called at the start of each test that modifies data
func resetMockDataUsers() {
	services.ResetOrderMockData()
	services.ResetUserMockData()
}

func TestGetUserWithOrders(t *testing.T) {
	resetMockDataUsers()
	
	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "Valid user ID returns user with orders",
			userID:         "750e8400-e29b-41d4-a716-446655440000",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var userOrders models.UserOrders
				if err := json.NewDecoder(w.Body).Decode(&userOrders); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if userOrders.User.ID != "750e8400-e29b-41d4-a716-446655440000" {
					t.Errorf("Expected user ID 750e8400-e29b-41d4-a716-446655440000, got %s", userOrders.User.ID)
				}
				if userOrders.User.Username != "johndoe" {
					t.Errorf("Expected username johndoe, got %s", userOrders.User.Username)
				}
				if len(userOrders.Orders) != 2 {
					t.Errorf("Expected 2 orders, got %d", len(userOrders.Orders))
				}
			},
		},
		{
			name:           "User with no or few orders returns orders array",
			userID:         "750e8400-e29b-41d4-a716-446655440002",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var userOrders models.UserOrders
				if err := json.NewDecoder(w.Body).Decode(&userOrders); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if userOrders.User.Username != "bobsmith" {
					t.Errorf("Expected username bobsmith, got %s", userOrders.User.Username)
				}
				// Note: May have orders from other tests, so just check it returns an array
				if userOrders.Orders == nil {
					t.Error("Expected orders array (even if empty), got nil")
				}
			},
		},
		{
			name:           "Non-existent user ID returns 404",
			userID:         "750e8400-e29b-41d4-a716-446655440099",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var errorResp models.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if errorResp.Code != "USER_NOT_FOUND" {
					t.Errorf("Expected error code USER_NOT_FOUND, got %s", errorResp.Code)
				}
			},
		},
		{
			name:           "Invalid UUID format returns 400",
			userID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var errorResp models.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if errorResp.Code != "INVALID_USER_ID" {
					t.Errorf("Expected error code INVALID_USER_ID, got %s", errorResp.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/user/"+tt.userID, nil)
			w := httptest.NewRecorder()

			GetUserWithOrders(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestGetUserLoyaltyPoints(t *testing.T) {
	resetMockDataUsers()
	
	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		expectedPoints int
	}{
		{
			name:           "valid user with loyalty points",
			userID:         "750e8400-e29b-41d4-a716-446655440000",
			expectedStatus: http.StatusOK,
			expectedPoints: 1500,
		},
		{
			name:           "valid user 2 with different points",
			userID:         "750e8400-e29b-41d4-a716-446655440001",
			expectedStatus: http.StatusOK,
			expectedPoints: 2300,
		},
		{
			name:           "user not found",
			userID:         "850e8400-e29b-41d4-a716-446655440099",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid UUID format",
			userID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/user/"+tt.userID+"/points", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetUserLoyaltyPoints)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]int
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not parse response: %v", err)
				}

				if points, ok := response["loyaltyPoints"]; !ok || points != tt.expectedPoints {
					t.Errorf("handler returned wrong loyalty points: got %v want %v", points, tt.expectedPoints)
				}
			}
		})
	}
}

func TestRedeemUserLoyaltyPoints(t *testing.T) {
	resetMockDataUsers()
	
	tests := []struct {
		name               string
		userID             string
		pointsToRedeem     int
		expectedStatus     int
		expectedRemaining  int
		skipBodyValidation bool
	}{
		{
			name:              "valid redemption",
			userID:            "750e8400-e29b-41d4-a716-446655440002",
			pointsToRedeem:    100,
			expectedStatus:    http.StatusOK,
			expectedRemaining: 400, // 500 - 100
		},
		{
			name:              "redeem all points",
			userID:            "750e8400-e29b-41d4-a716-446655440002",
			pointsToRedeem:    400, // After previous test
			expectedStatus:    http.StatusOK,
			expectedRemaining: 0,
		},
		{
			name:               "insufficient points",
			userID:             "750e8400-e29b-41d4-a716-446655440002",
			pointsToRedeem:     100,
			expectedStatus:     http.StatusBadRequest,
			skipBodyValidation: true,
		},
		{
			name:               "user not found",
			userID:             "850e8400-e29b-41d4-a716-446655440099",
			pointsToRedeem:     100,
			expectedStatus:     http.StatusNotFound,
			skipBodyValidation: true,
		},
		{
			name:               "invalid UUID format",
			userID:             "invalid-uuid",
			pointsToRedeem:     100,
			expectedStatus:     http.StatusBadRequest,
			skipBodyValidation: true,
		},
		{
			name:               "invalid points - zero",
			userID:             "750e8400-e29b-41d4-a716-446655440000",
			pointsToRedeem:     0,
			expectedStatus:     http.StatusBadRequest,
			skipBodyValidation: true,
		},
		{
			name:               "invalid points - negative",
			userID:             "750e8400-e29b-41d4-a716-446655440000",
			pointsToRedeem:     -50,
			expectedStatus:     http.StatusBadRequest,
			skipBodyValidation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := map[string]int{
				"points": tt.pointsToRedeem,
			}
			bodyBytes, _ := json.Marshal(requestBody)

			req, err := http.NewRequest("POST", "/user/"+tt.userID+"/points", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(RedeemUserLoyaltyPoints)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK && !tt.skipBodyValidation {
				var response map[string]int
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("could not parse response: %v", err)
				}

				if remaining, ok := response["remainingPoints"]; !ok || remaining != tt.expectedRemaining {
					t.Errorf("handler returned wrong remaining points: got %v want %v", remaining, tt.expectedRemaining)
				}
			}
		})
	}
}

func TestRedeemUserLoyaltyPoints_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/user/750e8400-e29b-41d4-a716-446655440000/points", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RedeemUserLoyaltyPoints)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCreateUser(t *testing.T) {
	resetMockDataUsers()
	
	tests := []struct {
		name           string
		requestBody    map[string]string
		authHeader     string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "valid user creation",
			requestBody: map[string]string{
				"username":  "newuser",
				"email":     "newuser@example.com",
				"firstname": "New",
				"lastname":  "User",
			},
			authHeader:     "Bearer valid_test_token_123456",
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var user map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &user)
				if err != nil {
					t.Fatalf("could not parse response: %v", err)
				}
				if user["username"] != "newuser" {
					t.Errorf("expected username 'newuser', got %v", user["username"])
				}
				if user["email"] != "newuser@example.com" {
					t.Errorf("expected email 'newuser@example.com', got %v", user["email"])
				}
				if user["id"] == nil || user["id"] == "" {
					t.Error("expected non-empty user ID")
				}
			},
		},
		{
			name: "missing username",
			requestBody: map[string]string{
				"email":     "test@example.com",
				"firstname": "Test",
				"lastname":  "User",
			},
			authHeader:     "Bearer valid_test_token_123456",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing email",
			requestBody: map[string]string{
				"username":  "testuser",
				"firstname": "Test",
				"lastname":  "User",
			},
			authHeader:     "Bearer valid_test_token_123456",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing firstname",
			requestBody: map[string]string{
				"username": "testuser",
				"email":    "test@example.com",
				"lastname": "User",
			},
			authHeader:     "Bearer valid_test_token_123456",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing lastname",
			requestBody: map[string]string{
				"username":  "testuser",
				"email":     "test@example.com",
				"firstname": "Test",
			},
			authHeader:     "Bearer valid_test_token_123456",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "username too short",
			requestBody: map[string]string{
				"username":  "ab",
				"email":     "test@example.com",
				"firstname": "Test",
				"lastname":  "User",
			},
			authHeader:     "Bearer valid_test_token_123456",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "username too long",
			requestBody: map[string]string{
				"username":  "thisusernameiswaytoolongandexceedsthirtychars",
				"email":     "test@example.com",
				"firstname": "Test",
				"lastname":  "User",
			},
			authHeader:     "Bearer valid_test_token_123456",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email format",
			requestBody: map[string]string{
				"username":  "testuser",
				"email":     "notanemail",
				"firstname": "Test",
				"lastname":  "User",
			},
			authHeader:     "Bearer valid_test_token_123456",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate username",
			requestBody: map[string]string{
				"username":  "johndoe",
				"email":     "different@example.com",
				"firstname": "John",
				"lastname":  "Different",
			},
			authHeader:     "Bearer valid_test_token_123456",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(CreateUser)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr)
			}
		})
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_test_token_123456")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestDeleteUser(t *testing.T) {
	resetMockDataUsers()
	
	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "delete existing user",
			userID:         "750e8400-e29b-41d4-a716-446655440002", // bobsmith
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete non-existent user returns 404",
			userID:         "750e8400-e29b-41d4-a716-446655440099",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid UUID format returns 400",
			userID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty user ID returns 400",
			userID:         "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/user/"+tt.userID, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(DeleteUser)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v. Body: %s", status, tt.expectedStatus, rr.Body.String())
			}

			// For 204 No Content, verify no body is returned
			if tt.expectedStatus == http.StatusNoContent && rr.Body.Len() > 0 {
				t.Errorf("handler returned body for 204 No Content response")
			}
		})
	}
}

func TestDeleteUser_CancelsPendingOrders(t *testing.T) {
	resetMockDataUsers()
	
	// This test verifies that deleting a user cancels all their PENDING orders
	// User 750e8400-e29b-41d4-a716-446655440000 (johndoe) has order 650e8400-e29b-41d4-a716-446655440000 which is PENDING
	
	// First, verify the order is PENDING
	orderReq, _ := http.NewRequest("GET", "/orders/650e8400-e29b-41d4-a716-446655440000", nil)
	orderRR := httptest.NewRecorder()
	orderHandler := http.HandlerFunc(GetOrderByID)
	orderHandler.ServeHTTP(orderRR, orderReq)
	
	if orderRR.Code != http.StatusOK {
		t.Fatalf("failed to get order: status %v", orderRR.Code)
	}
	
	var orderBefore map[string]interface{}
	json.Unmarshal(orderRR.Body.Bytes(), &orderBefore)
	if orderBefore["status"] != "PENDING" {
		t.Fatalf("expected order to be PENDING before deletion, got %v", orderBefore["status"])
	}
	
	// Delete the user
	deleteReq, _ := http.NewRequest("DELETE", "/user/750e8400-e29b-41d4-a716-446655440000", nil)
	deleteRR := httptest.NewRecorder()
	deleteHandler := http.HandlerFunc(DeleteUser)
	deleteHandler.ServeHTTP(deleteRR, deleteReq)
	
	if deleteRR.Code != http.StatusNoContent {
		t.Fatalf("failed to delete user: status %v, body: %s", deleteRR.Code, deleteRR.Body.String())
	}
	
	// Verify the order is now CANCELED
	orderReq2, _ := http.NewRequest("GET", "/orders/650e8400-e29b-41d4-a716-446655440000", nil)
	orderRR2 := httptest.NewRecorder()
	orderHandler.ServeHTTP(orderRR2, orderReq2)
	
	if orderRR2.Code != http.StatusOK {
		t.Fatalf("failed to get order after deletion: status %v", orderRR2.Code)
	}
	
	var orderAfter map[string]interface{}
	json.Unmarshal(orderRR2.Body.Bytes(), &orderAfter)
	if orderAfter["status"] != "CANCELED" {
		t.Errorf("expected order to be CANCELED after user deletion, got %v", orderAfter["status"])
	}
}

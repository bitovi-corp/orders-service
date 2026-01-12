package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	authmiddleware "github.com/bitovi-corp/auth-middleware-go/middleware"
	"github.com/Bitovi/example-go-server/internal/handlers"
	"github.com/Bitovi/example-go-server/internal/middleware"
	"github.com/Bitovi/example-go-server/internal/services"
)

func TestMain(m *testing.M) {
	// Reset mock data before running tests
	services.ResetOrderMockData()
	services.ResetUserMockData()
	
	// Run tests
	code := m.Run()
	
	os.Exit(code)
}

// TestOrderWorkflow implements the complete order workflow integration test
// as specified in order_workflow_test.md
func TestOrderWorkflow(t *testing.T) {
	// Reset mock data at the start of the test
	services.ResetOrderMockData()
	services.ResetUserMockData()
	
	// Helper function to make authenticated requests
	makeRequest := func(method, path string, body interface{}) *httptest.ResponseRecorder {
		var reqBody io.Reader
		if body != nil {
			jsonBytes, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}
			reqBody = bytes.NewReader(jsonBytes)
		}

		req := httptest.NewRequest(method, path, reqBody)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test_token_for_integration_test_12345")

		rr := httptest.NewRecorder()

		// Route to appropriate handler with middleware
		switch {
		case path == "/user":
			middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.CreateUser))(rr, req)
		case method == "POST" && len(path) > 7 && path[:7] == "/orders" && path[len(path)-7:] == "/submit":
			middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.CancelOrSubmitOrder))(rr, req)
		case method == "POST" && path == "/orders":
			middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.CreateOrder))(rr, req)
		case method == "GET" && len(path) > 8 && path[:8] == "/orders/":
			middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.GetOrderByID))(rr, req)
		case method == "PATCH" && len(path) > 8 && path[:8] == "/orders/":
			middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.UpdateOrder))(rr, req)
		case len(path) > 7 && path[len(path)-7:] == "/points":
			middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.GetUserLoyaltyPoints))(rr, req)
		case method == "DELETE" && len(path) > 6 && path[:6] == "/user/":
			middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.DeleteUser))(rr, req)
		case method == "GET" && len(path) > 6 && path[:6] == "/user/":
			middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.GetUserWithOrders))(rr, req)
		default:
			t.Fatalf("No handler found for %s %s", method, path)
		}

		return rr
	}

	// Step 1: Create a new user
	t.Log("Step 1: Create a new user")
	createUserBody := map[string]string{
		"email":     "test@example.com",
		"username":  "Test_User",
		"firstname": "Jane",
		"lastname":  "Doe",
	}
	resp := makeRequest("POST", "/user", createUserBody)
	if resp.Code != http.StatusCreated {
		t.Fatalf("Step 1 failed: Expected 201, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var user map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &user); err != nil {
		t.Fatalf("Step 1: Failed to parse user response: %v", err)
	}
	userID, ok := user["id"].(string)
	if !ok || userID == "" {
		t.Fatalf("Step 1: No user ID in response")
	}
	t.Logf("Created user with ID: %s", userID)

	// Step 2: Create a new order for the new user
	t.Log("Step 2: Create a new order")
	createOrderBody := map[string]interface{}{
		"userId": userID,
		"products": []map[string]interface{}{
			{"productId": "550e8400-e29b-41d4-a716-446655440000", "quantity": 1}, // Laptop
			{"productId": "550e8400-e29b-41d4-a716-446655440003", "quantity": 3}, // Notebook
		},
	}
	resp = makeRequest("POST", "/orders", createOrderBody)
	if resp.Code != http.StatusCreated {
		t.Fatalf("Step 2 failed: Expected 201, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var order map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &order); err != nil {
		t.Fatalf("Step 2: Failed to parse order response: %v", err)
	}
	orderID, ok := order["id"].(string)
	if !ok || orderID == "" {
		t.Fatalf("Step 2: No order ID in response")
	}
	t.Logf("Created order with ID: %s", orderID)

	// Step 3: Add 1 wireless mouse to the order
	t.Log("Step 3: Add wireless mouse to order")
	updateOrderBody := map[string]interface{}{
		"products": []map[string]interface{}{
			{"productId": "550e8400-e29b-41d4-a716-446655440001", "quantity": 1}, // Wireless Mouse
		},
	}
	resp = makeRequest("PATCH", "/orders/"+orderID, updateOrderBody)
	if resp.Code != http.StatusOK {
		t.Fatalf("Step 3 failed: Expected 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	if err := json.Unmarshal(resp.Body.Bytes(), &order); err != nil {
		t.Fatalf("Step 3: Failed to parse order response: %v", err)
	}
	products, ok := order["products"].([]interface{})
	if !ok {
		t.Fatalf("Step 3: No products array in response")
	}
	if len(products) != 3 {
		t.Fatalf("Step 3: Expected 3 products, got %d", len(products))
	}
	t.Logf("Order now has %d products", len(products))

	// Step 4: Check loyalty points (should be 0 before order submission)
	t.Log("Step 4: Check loyalty points before submission")
	resp = makeRequest("GET", "/user/"+userID+"/points", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("Step 4 failed: Expected 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var pointsResp map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &pointsResp); err != nil {
		t.Fatalf("Step 4: Failed to parse points response: %v", err)
	}
	loyaltyPoints, ok := pointsResp["loyaltyPoints"].(float64)
	if !ok {
		t.Fatalf("Step 4: No loyaltyPoints in response")
	}
	if loyaltyPoints != 0 {
		t.Fatalf("Step 4: Expected 0 loyalty points, got %.0f", loyaltyPoints)
	}
	t.Log("Loyalty points correctly at 0 before submission")

	// Step 5: Submit the order
	t.Log("Step 5: Submit the order")
	submitBody := map[string]string{
		"action": "SUBMIT",
	}
	resp = makeRequest("POST", "/orders/"+orderID+"/submit", submitBody)
	if resp.Code != http.StatusOK {
		t.Fatalf("Step 5 failed: Expected 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	if err := json.Unmarshal(resp.Body.Bytes(), &order); err != nil {
		t.Fatalf("Step 5: Failed to parse order response: %v", err)
	}
	status, ok := order["status"].(string)
	if !ok || status != "PROCESSING" {
		t.Fatalf("Step 5: Expected status PROCESSING, got %v", status)
	}
	t.Log("Order status changed to PROCESSING")

	// Step 6: Check the status of the order
	t.Log("Step 6: Verify order status")
	resp = makeRequest("GET", "/orders/"+orderID, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("Step 6 failed: Expected 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	if err := json.Unmarshal(resp.Body.Bytes(), &order); err != nil {
		t.Fatalf("Step 6: Failed to parse order response: %v", err)
	}
	status, ok = order["status"].(string)
	if !ok || status != "PROCESSING" {
		t.Fatalf("Step 6: Expected status PROCESSING, got %v", status)
	}
	totalPrice, ok := order["totalPrice"].(float64)
	if !ok {
		t.Fatalf("Step 6: No totalPrice in response")
	}
	t.Logf("Order status: %s, Total price: $%.2f", status, totalPrice)

	// Step 7: Check the loyalty points after submission
	t.Log("Step 7: Check loyalty points after submission")
	resp = makeRequest("GET", "/user/"+userID+"/points", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("Step 7 failed: Expected 200, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	if err := json.Unmarshal(resp.Body.Bytes(), &pointsResp); err != nil {
		t.Fatalf("Step 7: Failed to parse points response: %v", err)
	}
	loyaltyPoints, ok = pointsResp["loyaltyPoints"].(float64)
	if !ok {
		t.Fatalf("Step 7: No loyaltyPoints in response")
	}

	// Expected: Laptop ($1299.99) + Notebook ($19.99 × 3 = $59.97) + Mouse ($29.99) = $1389.95
	// Loyalty points: floor(1389.95 / 10) = 138
	expectedPoints := 138.0
	if loyaltyPoints != expectedPoints {
		t.Fatalf("Step 7: Expected %v loyalty points, got %.0f", expectedPoints, loyaltyPoints)
	}
	t.Logf("Loyalty points correctly awarded: %.0f points", loyaltyPoints)

	// Step 8: Create a second order (PENDING, not submitted)
	t.Log("Step 8: Create a second order")
	createOrder2Body := map[string]interface{}{
		"userId": userID,
		"products": []map[string]interface{}{
			{"productId": "550e8400-e29b-41d4-a716-446655440000", "quantity": 1}, // Laptop
			{"productId": "550e8400-e29b-41d4-a716-446655440003", "quantity": 3}, // Notebook
		},
	}
	resp = makeRequest("POST", "/orders", createOrder2Body)
	if resp.Code != http.StatusCreated {
		t.Fatalf("Step 8 failed: Expected 201, got %d. Body: %s", resp.Code, resp.Body.String())
	}

	var order2 map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &order2); err != nil {
		t.Fatalf("Step 8: Failed to parse order response: %v", err)
	}
	order2ID, ok := order2["id"].(string)
	if !ok || order2ID == "" {
		t.Fatalf("Step 8: No order ID in response")
	}
	t.Logf("Created second order with ID: %s", order2ID)

	// Cleanup: Delete test user
	t.Log("Cleanup: Deleting test user")
	resp = makeRequest("DELETE", "/user/"+userID, nil)
	if resp.Code != http.StatusNoContent {
		t.Fatalf("Cleanup failed: Expected 204, got %d. Body: %s", resp.Code, resp.Body.String())
	}
	t.Log("User deleted successfully")

	// Verify first order (PROCESSING) is still PROCESSING
	t.Log("Cleanup verification: Check first order status")
	resp = makeRequest("GET", "/orders/"+orderID, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("Cleanup verification failed: Expected 200, got %d", resp.Code)
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &order); err != nil {
		t.Fatalf("Cleanup verification: Failed to parse order response: %v", err)
	}
	status, _ = order["status"].(string)
	if status != "PROCESSING" {
		t.Fatalf("Cleanup verification: Expected first order to remain PROCESSING, got %s", status)
	}
	t.Log("First order (submitted) correctly remains PROCESSING")

	// Verify second order (PENDING) was CANCELED
	t.Log("Cleanup verification: Check second order status")
	resp = makeRequest("GET", "/orders/"+order2ID, nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("Cleanup verification failed: Expected 200, got %d", resp.Code)
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &order2); err != nil {
		t.Fatalf("Cleanup verification: Failed to parse order2 response: %v", err)
	}
	status2, _ := order2["status"].(string)
	if status2 != "CANCELED" {
		t.Fatalf("Cleanup verification: Expected second order to be CANCELED, got %s", status2)
	}
	t.Log("Second order (pending) correctly CANCELED on user deletion")

	t.Log("✅ Integration test completed successfully!")
}

func TestOrderWorkflowPlaceholder(t *testing.T) {
	t.Skip("Integration tests not yet implemented - see specification")
}

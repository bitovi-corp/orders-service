package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bitovi/example-go-server/internal/models"
)

func TestListProducts(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "GET request without limit returns default 20 products",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ProductListResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response.Limit != 20 {
					t.Errorf("Expected limit 20, got %d", response.Limit)
				}
				if response.Total != 5 {
					t.Errorf("Expected total 5, got %d", response.Total)
				}
				if len(response.Products) != 5 {
					t.Errorf("Expected 5 products, got %d", len(response.Products))
				}
			},
		},
		{
			name:           "GET request with valid limit",
			queryParams:    "?limit=2",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.ProductListResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response.Limit != 2 {
					t.Errorf("Expected limit 2, got %d", response.Limit)
				}
				if len(response.Products) != 2 {
					t.Errorf("Expected 2 products, got %d", len(response.Products))
				}
			},
		},
		{
			name:           "GET request with invalid limit returns 400",
			queryParams:    "?limit=150",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var errorResp models.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if errorResp.Code != "INVALID_LIMIT" {
					t.Errorf("Expected error code INVALID_LIMIT, got %s", errorResp.Code)
				}
			},
		},
		{
			name:           "POST request returns 405",
			queryParams:    "",
			expectedStatus: http.StatusMethodNotAllowed,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := http.MethodGet
			if tt.expectedStatus == http.StatusMethodNotAllowed {
				method = http.MethodPost
			}

			req := httptest.NewRequest(method, "/products"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			ListProducts(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestGetProductByID(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "Valid product ID returns product",
			productID:      "550e8400-e29b-41d4-a716-446655440000",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var product models.Product
				if err := json.NewDecoder(w.Body).Decode(&product); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if product.ID != "550e8400-e29b-41d4-a716-446655440000" {
					t.Errorf("Expected product ID 550e8400-e29b-41d4-a716-446655440000, got %s", product.ID)
				}
				if product.Name != "Laptop" {
					t.Errorf("Expected product name Laptop, got %s", product.Name)
				}
			},
		},
		{
			name:           "Non-existent product ID returns 404",
			productID:      "550e8400-e29b-41d4-a716-446655440099",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var errorResp models.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if errorResp.Code != "PRODUCT_NOT_FOUND" {
					t.Errorf("Expected error code PRODUCT_NOT_FOUND, got %s", errorResp.Code)
				}
			},
		},
		{
			name:           "Invalid UUID format returns 400",
			productID:      "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var errorResp models.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if errorResp.Code != "INVALID_PRODUCT_ID" {
					t.Errorf("Expected error code INVALID_PRODUCT_ID, got %s", errorResp.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.productID, nil)
			w := httptest.NewRecorder()

			GetProductByID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

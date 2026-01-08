package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bitovi/example-go-server/internal/models"
	"github.com/Bitovi/example-go-server/internal/services"
	"github.com/google/uuid"
)

var productService = services.NewProductService()

// ListProducts implements GET /products endpoint as defined in api/openapi.yaml
func ListProducts(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Parse limit query parameter
	limit := 20 // default value as per spec
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 1 || parsedLimit > 100 {
			writeErrorResponse(w, http.StatusBadRequest, "INVALID_LIMIT", "Limit must be between 1 and 100", "")
			return
		}
		limit = parsedLimit
	}

	// Get products from service
	products, total := productService.ListProducts(limit)

	// Prepare response
	response := models.ProductListResponse{
		Products: products,
		Total:    total,
		Limit:    limit,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding products list response: %v", err)
	}
}

// GetProductByID implements GET /products/{productId} endpoint as defined in api/openapi.yaml
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", "")
		return
	}

	// Extract product ID from URL path
	// URL format: /products/{productId}
	path := strings.TrimPrefix(r.URL.Path, "/products/")
	productID := strings.Split(path, "/")[0]

	if productID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_PRODUCT_ID", "Product ID is required", "")
		return
	}

	// Basic UUID format validation (simple check)
	if !isValidUUID(productID) {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_PRODUCT_ID", "Invalid product ID format", "Product ID must be a valid UUID")
		return
	}

	// Get product from service
	product, err := productService.GetProductByID(productID)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			writeErrorResponse(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "The requested product could not be found", "")
			return
		}
		log.Printf("Error retrieving product: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", "")
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("Error encoding product response: %v", err)
	}
}

// writeErrorResponse writes an error response according to the Error schema in openapi.yaml
func writeErrorResponse(w http.ResponseWriter, statusCode int, code, message, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	errorResp := models.ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
	
	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}

// isValidUUID performs UUID format validation using google/uuid
func isValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}

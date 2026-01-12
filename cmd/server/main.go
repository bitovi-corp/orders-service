package main

import (
	"log"
	"net/http"

	authmiddleware "github.com/bitovi-corp/auth-middleware-go/middleware"
	"github.com/Bitovi/example-go-server/internal/handlers"
	"github.com/Bitovi/example-go-server/internal/middleware"
)

func main() {
	// Register routes according to api/openapi.yaml
	// Health check endpoint - no auth required
	http.HandleFunc("/health", middleware.LoggingMiddleware(handlers.HealthCheck))
	
	// Product endpoints - auth required
	http.HandleFunc("/products", middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.ListProducts)))
	http.HandleFunc("/products/", middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.GetProductByID)))
	
	// Order endpoints - auth required
	http.HandleFunc("/orders/", middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handleOrdersWithID)))
	http.HandleFunc("/orders", middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handleOrders)))
	
	// User endpoints - auth required
	http.HandleFunc("/user", middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handlers.CreateUser)))
	http.HandleFunc("/user/", middleware.LoggingMiddleware(authmiddleware.AuthMiddleware(handleUserRoutes)))

	// Start server
	port := ":8080"
	log.Printf("Starting Example Server API v1.0.0 on port %s", port)
	log.Printf("Endpoints available:")
	log.Printf("  - GET http://localhost%s/health (no auth)", port)
	log.Printf("  - GET http://localhost%s/products (auth required)", port)
	log.Printf("  - GET http://localhost%s/products/{productId} (auth required)", port)
	log.Printf("  - GET http://localhost%s/orders (auth required)", port)
	log.Printf("  - POST http://localhost%s/orders (auth required)", port)
	log.Printf("  - GET http://localhost%s/orders/{orderId} (auth required)", port)
	log.Printf("  - PATCH http://localhost%s/orders/{orderId} (auth required)", port)
	log.Printf("  - POST http://localhost%s/orders/{orderId}/submit (auth required)", port)
	log.Printf("  - GET http://localhost%s/user/{userId} (auth required)", port)
	log.Printf("  - DELETE http://localhost%s/user/{userId} (auth required)", port)
	log.Printf("  - GET http://localhost%s/user/{userId}/points (auth required)", port)
	log.Printf("  - POST http://localhost%s/user/{userId}/points (auth required)", port)
	log.Printf("  - POST http://localhost%s/user (auth required)", port)
	log.Printf("")
	log.Printf("Authentication: Include 'Authorization: Bearer {token}' header")
	log.Printf("Global middlewares: Logging enabled for all requests")
	
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// handleOrders routes to the appropriate handler based on method
func handleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlers.ListOrders(w, r)
	case http.MethodPost:
		handlers.CreateOrder(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleOrdersWithID routes to the appropriate handler based on method for /orders/{orderId}
func handleOrdersWithID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	// Check if it's the submit endpoint: /orders/{orderId}/submit
	if len(path) > 7 && path[len(path)-7:] == "/submit" {
		if r.Method == http.MethodPost {
			handlers.CancelOrSubmitOrder(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	
	// Regular /orders/{orderId} endpoints
	switch r.Method {
	case http.MethodGet:
		handlers.GetOrderByID(w, r)
	case http.MethodPatch:
		handlers.UpdateOrder(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleUserRoutes routes user endpoints based on path and method
func handleUserRoutes(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	// Check if it's the points endpoint
	if len(path) > 7 && path[len(path)-7:] == "/points" {
		// /user/{userId}/points
		switch r.Method {
		case http.MethodGet:
			handlers.GetUserLoyaltyPoints(w, r)
		case http.MethodPost:
			handlers.RedeemUserLoyaltyPoints(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		// /user/{userId}
		switch r.Method {
		case http.MethodGet:
			handlers.GetUserWithOrders(w, r)
		case http.MethodDelete:
			handlers.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
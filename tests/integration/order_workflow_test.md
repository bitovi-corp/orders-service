# Order Workflow Integration Test Specification

## Test: Complete Order Flow with Loyalty Points

### Description
This test validates the complete order creation workflow including user creation, order management, and loyalty points calculation.

### Prerequisites
- Server running on localhost:8080
- Database seeded with product catalog
- Authentication system configured

### Test Steps

#### 1. Create a new user
- **Action**: POST /user
- **Request Body**:
  ```json
  {
    "email": "test@example.com",
    "username": "Test_User",
	"firstname": "Jane",
	"lastname": "Doe"
  }
  ```
- **Expected Response**: 201 Created
- **Capture**: `userId` from response

#### 2. Create a new order for the new user
- **Action**: POST /orders
- **Request Body**:
  ```json
  {
    "userId": "{userId from step 1}",
    "products": [
      {
        "productId": "550e8400-e29b-41d4-a716-446655440000",
        "quantity": 1
      },
      {
        "productId": "550e8400-e29b-41d4-a716-446655440003",
        "quantity": 3
      }
    ]
  }
  ```
- **Expected Response**: 201 Created
- **Capture**: `orderId` from response

#### 3. Add 1 wireless mouse to the order
- **Action**: PATCH /orders/{orderId}
- **Request Body**:
  ```json
  {
    "products": [
      {"productId": "550e8400-e29b-41d4-a716-446655440001", "quantity": 1}
    ]
  }
  ```
- **Expected Response**: 200 OK
- **Expected**: Order should now contain 3 distinct items

#### 4. Check loyalty points (should be 0 before order submission)
- **Action**: GET /user/{userId}/points
- **Expected Response**: 200 OK
- **Expected Body**:
  ```json
  {
    "loyaltyPoints": 0
  }
  ```

#### 5. Submit the order
- **Action**: POST /orders/{orderId}/submit
- **Request Body**:
  ```json
  {
    "action": "SUBMIT"
  }
  ```
- **Expected Response**: 200 OK
- **Expected**: Order status changes to "PROCESSING"

#### 6. Check the status of the order
- **Action**: GET /orders/{orderId}
- **Expected Response**: 200 OK
- **Expected Body**:
  ```json
  {
    "id": "{order ID}",
    "status": "PROCESSING",
    "products": [...],
    "totalPrice": {calculated amount}
  }
  ```
- **Validation**: `status` field must equal "PROCESSING"

#### 7. Check the loyalty points for the order
- **Action**: GET /user/{userId}/points
- **Expected Response**: 200 OK
- **Calculation Logic**:
  - Laptop: $1299.99 × 1 = $1299.99
  - Notebook: $19.99 × 3 = $59.97
  - Wireless Mouse: $29.99 × 1 = $29.99
  - **Total**: $1,389.95
  - **Loyalty Points**: 1 point per $10 spent = 138 points (rounded down)
- **Expected Body**:
  ```json
  {
    "loyaltyPoints": 138
  }
  ```

#### 8. Create a new order for the new user
- **Action**: POST /orders
- **Request Body**:
  ```json
  {
    "userId": "{userId from step 1}",
    "products": [
      {
        "productId": "550e8400-e29b-41d4-a716-446655440000",
        "quantity": 1
      },
      {
        "productId": "550e8400-e29b-41d4-a716-446655440003",
        "quantity": 3
      }
    ]
  }
  ```
- **Expected Response**: 201 Created
- **Capture**: `orderId` from response

### Cleanup
- Delete test user: DELETE /user/{userId}
- Verify the first order that was submitted is still being processed, verify he second order that was submitted is cancelled

### Notes
- This test requires authentication token (use test auth middleware)
- Product prices must match the calculation above
- Loyalty points calculation: floor(total_amount / 10)
- Test should be idempotent and clean up after itself
# E-Commerce API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
The API uses JWT (JSON Web Token) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## API Endpoints

### Authentication

#### Register User
```http
POST /auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "role": "customer",
    "created_at": "2024-01-10T10:00:00Z",
    "updated_at": "2024-01-10T10:00:00Z"
  }
}
```

#### Login
```http
POST /auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "customer",
    "created_at": "2024-01-10T10:00:00Z",
    "updated_at": "2024-01-10T10:00:00Z"
  }
}
```

#### Get Profile
```http
GET /auth/profile
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "role": "customer",
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T10:00:00Z"
}
```

### Products

#### List Products
```http
GET /products
```

**Query Parameters:**
- `page` (int, optional): Page number (default: 1)
- `limit` (int, optional): Items per page (default: 20)
- `category_id` (int, optional): Filter by category
- `min_price` (float, optional): Minimum price filter
- `max_price` (float, optional): Maximum price filter
- `search` (string, optional): Search in name and description

**Response:**
```json
{
  "products": [
    {
      "id": 1,
      "name": "Laptop Pro",
      "description": "High-performance laptop",
      "price": 1299.99,
      "stock_quantity": 50,
      "category_id": 1,
      "sku": "LP-001",
      "image_url": "https://example.com/laptop.jpg",
      "is_active": true,
      "created_at": "2024-01-10T10:00:00Z",
      "updated_at": "2024-01-10T10:00:00Z"
    }
  ],
  "count": 1
}
```

#### Get Product
```http
GET /products/{id}
```

**Response:**
```json
{
  "id": 1,
  "name": "Laptop Pro",
  "description": "High-performance laptop",
  "price": 1299.99,
  "stock_quantity": 50,
  "category_id": 1,
  "sku": "LP-001",
  "image_url": "https://example.com/laptop.jpg",
  "is_active": true,
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T10:00:00Z"
}
```

#### Create Product (Admin Only)
```http
POST /products
```

**Headers:**
```
Authorization: Bearer <admin-token>
```

**Request Body:**
```json
{
  "name": "New Product",
  "description": "Product description",
  "price": 99.99,
  "stock_quantity": 100,
  "category_id": 1,
  "sku": "NP-001",
  "image_url": "https://example.com/product.jpg"
}
```

#### Update Product (Admin Only)
```http
PUT /products/{id}
```

#### Delete Product (Admin Only)
```http
DELETE /products/{id}
```

### Shopping Cart

#### Get Cart
```http
GET /cart
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "items": [
    {
      "id": 1,
      "user_id": 1,
      "product_id": 1,
      "quantity": 2,
      "created_at": "2024-01-10T10:00:00Z",
      "updated_at": "2024-01-10T10:00:00Z",
      "product": {
        "id": 1,
        "name": "Laptop Pro",
        "price": 1299.99
      }
    }
  ],
  "total_items": 2,
  "total_price": 2599.98
}
```

#### Add to Cart
```http
POST /cart
```

**Request Body:**
```json
{
  "product_id": 1,
  "quantity": 2
}
```

#### Update Cart Item
```http
PUT /cart/{product_id}
```

**Request Body:**
```json
{
  "quantity": 3
}
```

#### Remove from Cart
```http
DELETE /cart/{product_id}
```

#### Clear Cart
```http
DELETE /cart
```

### Orders

#### Create Order
```http
POST /orders
```

**Request Body:**
```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    }
  ],
  "shipping_address": "123 Main St, City, State 12345",
  "billing_address": "123 Main St, City, State 12345"
}
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "order_number": "ORD-abc12345",
  "status": "pending",
  "total_amount": 2599.98,
  "shipping_address": "123 Main St, City, State 12345",
  "billing_address": "123 Main St, City, State 12345",
  "created_at": "2024-01-10T10:00:00Z",
  "updated_at": "2024-01-10T10:00:00Z",
  "order_items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": 1,
      "quantity": 2,
      "unit_price": 1299.99,
      "total_price": 2599.98
    }
  ]
}
```

#### Get User Orders
```http
GET /orders
```

#### Get Order
```http
GET /orders/{id}
```

#### Process Payment
```http
POST /payments
```

**Request Body:**
```json
{
  "order_id": 1,
  "payment_method": "credit_card"
}
```

### Reviews

#### Get Product Reviews
```http
GET /products/{product_id}/reviews
```

**Response:**
```json
{
  "reviews": [
    {
      "id": 1,
      "product_id": 1,
      "user_id": 1,
      "rating": 5,
      "comment": "Excellent product!",
      "created_at": "2024-01-10T10:00:00Z",
      "updated_at": "2024-01-10T10:00:00Z",
      "user": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe"
      }
    }
  ],
  "count": 1
}
```

#### Create Review
```http
POST /products/{product_id}/reviews
```

**Request Body:**
```json
{
  "rating": 5,
  "comment": "Excellent product!"
}
```

#### Update Review
```http
PUT /reviews/{id}
```

#### Delete Review
```http
DELETE /reviews/{id}
```

### Admin Endpoints

#### Get All Orders
```http
GET /admin/orders
```

#### Update Order Status
```http
PUT /admin/orders/{id}
```

**Request Body:**
```json
{
  "status": "shipped"
}
```

### Health Checks

#### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-10T10:00:00Z",
  "version": "1.0.0",
  "checks": {
    "database": {
      "status": "healthy",
      "message": "Database connection successful",
      "latency": "5ms"
    }
  },
  "uptime": "2h30m45s"
}
```

#### Readiness Check
```http
GET /ready
```

#### Liveness Check
```http
GET /live
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message description"
}
```

**HTTP Status Codes:**
- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `429` - Too Many Requests
- `500` - Internal Server Error
- `503` - Service Unavailable

## Rate Limiting

API requests are rate-limited to 100 requests per minute per IP address.

## CORS

The API supports Cross-Origin Resource Sharing (CORS) with appropriate headers.

## Security

- All passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- HTTPS is recommended in production
- Security headers are included in all responses

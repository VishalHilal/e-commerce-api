# API Examples

This document provides examples of how to use the e-commerce API endpoints.

## Authentication

### Register a new user
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "customer",
    "created_at": "2024-01-10T10:00:00Z"
  }
}
```

## Products

### List all products
```bash
curl "http://localhost:8080/products?page=1&limit=10&search=laptop&min_price=500&max_price=2000"
```

### Get product details
```bash
curl http://localhost:8080/products/1
```

### Create a product (admin only)
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Gaming Laptop",
    "description": "High-performance gaming laptop",
    "price": 1499.99,
    "stock_quantity": 25,
    "category_id": 1,
    "sku": "GL-001",
    "image_url": "https://example.com/laptop.jpg"
  }'
```

## Shopping Cart

### Get cart
```bash
curl -X GET http://localhost:8080/cart \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Add item to cart
```bash
curl -X POST http://localhost:8080/cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

### Update cart item
```bash
curl -X PUT http://localhost:8080/cart/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "quantity": 3
  }'
```

### Remove item from cart
```bash
curl -X DELETE http://localhost:8080/cart/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Orders

### Create order
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      },
      {
        "product_id": 2,
        "quantity": 1
      }
    ],
    "shipping_address": "123 Main St, City, State 12345",
    "billing_address": "123 Main St, City, State 12345"
  }'
```

### Get user orders
```bash
curl -X GET http://localhost:8080/orders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Process payment
```bash
curl -X POST http://localhost:8080/payments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "order_id": 1,
    "payment_method": "credit_card"
  }'
```

## Reviews

### Get product reviews
```bash
curl http://localhost:8080/products/1/reviews
```

### Create review
```bash
curl -X POST http://localhost:8080/products/1/reviews \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "rating": 5,
    "comment": "Excellent product! Highly recommended."
  }'
```

### Update review
```bash
curl -X PUT http://localhost:8080/reviews/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "rating": 4,
    "comment": "Good product, but shipping was slow."
  }'
```

## Admin Operations

### Get all orders (admin)
```bash
curl -X GET http://localhost:8080/admin/orders \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

### Update order status (admin)
```bash
curl -X PUT http://localhost:8080/admin/orders/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -d '{
    "status": "shipped"
  }'
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Pagination

List endpoints support pagination:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20)

## Search and Filtering

Products endpoint supports:
- `search` - Search in name and description
- `category_id` - Filter by category
- `min_price` - Minimum price filter
- `max_price` - Maximum price filter

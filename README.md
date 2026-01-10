# E-Commerce API

A comprehensive e-commerce API built with Go, Chi router, and PostgreSQL.

## Features

### Authentication & Authorization
- JWT-based authentication
- Role-based access control (customer, admin)
- User registration and login
- Profile management

### Product Management
- CRUD operations for products
- Product search and filtering
- Category-based filtering
- Price range filtering
- Stock management

### Shopping Cart
- Add items to cart
- Update cart quantities
- Remove items from cart
- Clear cart
- Cart total calculation

### Order Management
- Create orders from cart items
- Order status tracking
- Order history for users
- Admin order management

### Payment Processing
- Payment creation and processing
- Payment status tracking
- Transaction management

### Admin Features
- Product management
- Order management
- User management
- Sales analytics

## API Endpoints

### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - User login
- `GET /auth/profile` - Get user profile
- `PUT /auth/profile` - Update user profile

### Products
- `GET /products` - List products (with search/filter)
- `GET /products/{id}` - Get product details
- `POST /products` - Create product (admin only)
- `PUT /products/{id}` - Update product (admin only)
- `DELETE /products/{id}` - Delete product (admin only)

### Cart
- `GET /cart` - Get user cart
- `POST /cart` - Add item to cart
- `PUT /cart/{product_id}` - Update cart item
- `DELETE /cart/{product_id}` - Remove item from cart
- `DELETE /cart` - Clear cart

### Orders
- `POST /orders` - Create order
- `GET /orders` - Get user orders
- `GET /orders/{id}` - Get order details
- `POST /payments` - Process payment

### Admin
- `GET /admin/orders` - Get all orders
- `PUT /admin/orders/{id}` - Update order status

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Set up PostgreSQL database and run migrations:
```bash
psql -d your_database -f migrations/001_initial_schema.sql
```

3. Configure environment variables:
```bash
export GOOSE_DBSTRING="host=localhost user=postgres password=postgres dbname=ecom sslmode=disable"
export JWT_SECRET="your-secret-key-here"
```

4. Run the application:
```bash
go run cmd/main.go
```

## Database Schema

The API uses the following main tables:
- `users` - User accounts and authentication
- `products` - Product catalog
- `categories` - Product categories
- `cart_items` - Shopping cart items
- `orders` - Customer orders
- `order_items` - Order line items
- `payments` - Payment records
- `product_reviews` - Product reviews and ratings

## Security

- Passwords are hashed using bcrypt
- JWT tokens for authentication
- Role-based access control
- Input validation and sanitization

## Dependencies

- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `golang.org/x/crypto` - Cryptographic functions
- `github.com/google/uuid` - UUID generation

## Development

The project follows a clean architecture pattern with:
- Domain models in `internal/models`
- Business logic in services (`internal/*`)
- HTTP handlers in `internal/*/handlers`
- Database adapters in `internal/adapters/postgresql`
- Authentication logic in `internal/auth`

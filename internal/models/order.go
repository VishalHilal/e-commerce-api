package models

import (
	"time"
)

type Order struct {
	ID              int         `json:"id"`
	UserID          int         `json:"user_id"`
	OrderNumber     string      `json:"order_number"`
	Status          string      `json:"status"`
	TotalAmount     float64     `json:"total_amount"`
	ShippingAddress string      `json:"shipping_address"`
	BillingAddress  string      `json:"billing_address"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	OrderItems      []OrderItem `json:"order_items,omitempty"`
	User            *User       `json:"user,omitempty"`
}

type OrderItem struct {
	ID         int      `json:"id"`
	OrderID    int      `json:"order_id"`
	ProductID  int      `json:"product_id"`
	Quantity   int      `json:"quantity"`
	UnitPrice  float64  `json:"unit_price"`
	TotalPrice float64  `json:"total_price"`
	Product    *Product `json:"product,omitempty"`
}

type CreateOrderRequest struct {
	Items           []OrderItemRequest `json:"items" validate:"required,min=1"`
	ShippingAddress string             `json:"shipping_address" validate:"required"`
	BillingAddress  string             `json:"billing_address" validate:"required"`
}

type OrderItemRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed shipped delivered cancelled"`
}

type Payment struct {
	ID            int       `json:"id"`
	OrderID       int       `json:"order_id"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"`
	Amount        float64   `json:"amount"`
	TransactionID string    `json:"transaction_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreatePaymentRequest struct {
	OrderID       int    `json:"order_id" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
}

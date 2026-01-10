package models

import (
	"time"
)

type Product struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CategoryID    int       `json:"category_id"`
	SKU           string    `json:"sku"`
	ImageURL      string    `json:"image_url,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description,omitempty"`
	Price         float64 `json:"price" validate:"required,gt=0"`
	StockQuantity int     `json:"stock_quantity" validate:"required,gte=0"`
	CategoryID    int     `json:"category_id" validate:"required"`
	SKU           string  `json:"sku" validate:"required"`
	ImageURL      string  `json:"image_url,omitempty"`
}

type UpdateProductRequest struct {
	Name          *string  `json:"name,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Price         *float64 `json:"price,omitempty"`
	StockQuantity *int     `json:"stock_quantity,omitempty"`
	CategoryID    *int     `json:"category_id,omitempty"`
	ImageURL      *string  `json:"image_url,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}

type ProductFilter struct {
	CategoryID *int     `json:"category_id,omitempty"`
	MinPrice   *float64 `json:"min_price,omitempty"`
	MaxPrice   *float64 `json:"max_price,omitempty"`
	IsActive   *bool    `json:"is_active,omitempty"`
	Search     string   `json:"search,omitempty"`
	Page       int      `json:"page,omitempty"`
	Limit      int      `json:"limit,omitempty"`
}

package cart

import (
	"context"
	"fmt"

	"github.com/VishalHilal/e-commerce-api/internal/models"
)

type Repository interface {
	AddToCart(ctx context.Context, userID int, req models.AddToCartRequest) (*models.CartItem, error)
	GetCartItems(ctx context.Context, userID int) ([]models.CartItem, error)
	UpdateCartItem(ctx context.Context, userID, productID int, quantity int) error
	RemoveFromCart(ctx context.Context, userID, productID int) error
	ClearCart(ctx context.Context, userID int) error
	GetProductByID(ctx context.Context, productID int) (*models.Product, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddToCart(ctx context.Context, userID int, req models.AddToCartRequest) (*models.CartItem, error) {
	product, err := s.repo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if !product.IsActive {
		return nil, fmt.Errorf("product is not available")
	}

	if product.StockQuantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	cartItem, err := s.repo.AddToCart(ctx, userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to add to cart: %w", err)
	}

	return cartItem, nil
}

func (s *Service) GetCart(ctx context.Context, userID int) (*models.CartResponse, error) {
	items, err := s.repo.GetCartItems(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	var totalItems int
	var totalPrice float64

	for _, item := range items {
		totalItems += item.Quantity
		totalPrice += float64(item.Quantity) * item.Product.Price
	}

	return &models.CartResponse{
		Items:      items,
		TotalItems: totalItems,
		TotalPrice: totalPrice,
	}, nil
}

func (s *Service) UpdateCartItem(ctx context.Context, userID, productID int, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if !product.IsActive {
		return fmt.Errorf("product is not available")
	}

	if product.StockQuantity < quantity {
		return fmt.Errorf("insufficient stock")
	}

	if err := s.repo.UpdateCartItem(ctx, userID, productID, quantity); err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	return nil
}

func (s *Service) RemoveFromCart(ctx context.Context, userID, productID int) error {
	if err := s.repo.RemoveFromCart(ctx, userID, productID); err != nil {
		return fmt.Errorf("failed to remove from cart: %w", err)
	}
	return nil
}

func (s *Service) ClearCart(ctx context.Context, userID int) error {
	if err := s.repo.ClearCart(ctx, userID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}

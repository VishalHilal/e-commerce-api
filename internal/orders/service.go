package orders

import (
	"context"
	"fmt"

	"github.com/VishalHilal/e-commerce-api/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateOrder(ctx context.Context, order models.CreateOrderRequest, userID int) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error)
	GetOrderByID(ctx context.Context, id int) (*models.Order, error)
	UpdateOrderStatus(ctx context.Context, id int, status string) error
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	CreatePayment(ctx context.Context, payment models.CreatePaymentRequest) (*models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, req models.CreateOrderRequest, userID int) (*models.Order, error) {
	orderNumber := "ORD-" + uuid.New().String()[:8]

	order, err := s.repo.CreateOrder(ctx, req, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	order.OrderNumber = orderNumber
	return order, nil
}

func (s *Service) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	orders, err := s.repo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	return orders, nil
}

func (s *Service) GetOrder(ctx context.Context, orderID int, userID int) (*models.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if order.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to order")
	}

	return order, nil
}

func (s *Service) UpdateOrderStatus(ctx context.Context, orderID int, status string) error {
	validStatuses := []string{"pending", "confirmed", "shipped", "delivered", "cancelled"}
	isValid := false
	for _, s := range validStatuses {
		if s == status {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid order status: %s", status)
	}

	if err := s.repo.UpdateOrderStatus(ctx, orderID, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (s *Service) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	orders, err := s.repo.GetAllOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}
	return orders, nil
}

func (s *Service) ProcessPayment(ctx context.Context, req models.CreatePaymentRequest) (*models.Payment, error) {
	payment, err := s.repo.CreatePayment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	payment.PaymentStatus = "completed"
	if err := s.repo.UpdatePaymentStatus(ctx, payment.ID, payment.PaymentStatus); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	if err := s.repo.UpdateOrderStatus(ctx, req.OrderID, "confirmed"); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return payment, nil
}

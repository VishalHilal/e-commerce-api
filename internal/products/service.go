package products

import (
	"context"
	"fmt"

	"github.com/VishalHilal/e-commerce-api/internal/models"
)

type Repository interface {
	CreateProduct(ctx context.Context, product models.CreateProductRequest) (*models.Product, error)
	GetProducts(ctx context.Context, filter models.ProductFilter) ([]models.Product, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	UpdateProduct(ctx context.Context, id int, product models.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id int) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListProducts(ctx context.Context, filter models.ProductFilter) ([]models.Product, error) {
	return s.repo.GetProducts(ctx, filter)
}

func (s *Service) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (s *Service) CreateProduct(ctx context.Context, req models.CreateProductRequest) (*models.Product, error) {
	product, err := s.repo.CreateProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return product, nil
}

func (s *Service) UpdateProduct(ctx context.Context, id int, req models.UpdateProductRequest) error {
	_, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if err := s.repo.UpdateProduct(ctx, id, req); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (s *Service) DeleteProduct(ctx context.Context, id int) error {
	_, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

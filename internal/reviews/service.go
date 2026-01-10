package reviews

import (
	"context"
	"fmt"

	"github.com/VishalHilal/e-commerce-api/internal/models"
)

type Repository interface {
	CreateReview(ctx context.Context, review models.CreateReviewRequest, userID int) (*models.ProductReview, error)
	GetProductReviews(ctx context.Context, productID int) ([]models.ProductReview, error)
	UpdateReview(ctx context.Context, reviewID, userID int, req models.UpdateReviewRequest) (*models.ProductReview, error)
	DeleteReview(ctx context.Context, reviewID, userID int) error
	GetUserReview(ctx context.Context, productID, userID int) (*models.ProductReview, error)
	GetProductAverageRating(ctx context.Context, productID int) (float64, int, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateReview(ctx context.Context, req models.CreateReviewRequest, userID int) (*models.ProductReview, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	existingReview, err := s.repo.GetUserReview(ctx, req.ProductID, userID)
	if err == nil && existingReview != nil {
		return nil, fmt.Errorf("user has already reviewed this product")
	}

	review, err := s.repo.CreateReview(ctx, req, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

func (s *Service) GetProductReviews(ctx context.Context, productID int) ([]models.ProductReview, error) {
	reviews, err := s.repo.GetProductReviews(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product reviews: %w", err)
	}
	return reviews, nil
}

func (s *Service) UpdateReview(ctx context.Context, reviewID, userID int, req models.UpdateReviewRequest) (*models.ProductReview, error) {
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	review, err := s.repo.UpdateReview(ctx, reviewID, userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	return review, nil
}

func (s *Service) DeleteReview(ctx context.Context, reviewID, userID int) error {
	if err := s.repo.DeleteReview(ctx, reviewID, userID); err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}

func (s *Service) GetProductWithReviews(ctx context.Context, productID int) (*models.ProductWithReviews, error) {
	reviews, err := s.repo.GetProductReviews(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product reviews: %w", err)
	}

	avgRating, reviewCount, err := s.repo.GetProductAverageRating(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product rating: %w", err)
	}

	return &models.ProductWithReviews{
		Reviews:       reviews,
		AverageRating: avgRating,
		ReviewCount:   reviewCount,
	}, nil
}

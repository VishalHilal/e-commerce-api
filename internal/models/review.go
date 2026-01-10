package models

import (
	"time"
)

type ProductReview struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	UserID    int       `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user,omitempty"`
}

type CreateReviewRequest struct {
	ProductID int    `json:"product_id" validate:"required"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Comment   string `json:"comment,omitempty"`
}

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Comment *string `json:"comment,omitempty"`
}

type ProductWithReviews struct {
	Product
	AverageRating float64         `json:"average_rating"`
	ReviewCount   int             `json:"review_count"`
	Reviews       []ProductReview `json:"reviews,omitempty"`
}

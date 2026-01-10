package users

import (
	"context"
	"fmt"

	"github.com/VishalHilal/e-commerce-api/internal/auth"
	"github.com/VishalHilal/e-commerce-api/internal/models"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.CreateUserRequest) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, id int, user models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type Service struct {
	repo   Repository
	jwtSvc *auth.JWTService
}

func NewService(repo Repository, jwtSvc *auth.JWTService) *Service {
	return &Service{
		repo:   repo,
		jwtSvc: jwtSvc,
	}
}

func (s *Service) Register(ctx context.Context, req models.CreateUserRequest) (*models.LoginResponse, error) {
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userReq := models.CreateUserRequest{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}

	user, err := s.repo.CreateUser(ctx, userReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.jwtSvc.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *Service) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := s.jwtSvc.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *Service) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID int, req models.User) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	return s.repo.UpdateUser(ctx, userID, *user)
}

package postgresql

import (
	"context"
	"fmt"

	"github.com/VishalHilal/e-commerce-api/internal/models"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, email, first_name, last_name, phone, role, created_at, updated_at
	`

	var user models.User
	err := r.db.QueryRow(ctx, query,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
		req.Phone,
		"customer",
	).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id int, user models.User) error {
	query := `
		UPDATE users
		SET first_name = $2, last_name = $3, phone = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id, user.FirstName, user.LastName, user.Phone)
	return err
}

func (r *Repository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *Repository) CreateProduct(ctx context.Context, req models.CreateProductRequest) (*models.Product, error) {
	query := `
		INSERT INTO products (name, description, price, stock_quantity, category_id, sku, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, description, price, stock_quantity, category_id, sku, image_url, is_active, created_at, updated_at
	`

	var product models.Product
	err := r.db.QueryRow(ctx, query,
		req.Name,
		req.Description,
		req.Price,
		req.StockQuantity,
		req.CategoryID,
		req.SKU,
		req.ImageURL,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.CategoryID,
		&product.SKU,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *Repository) GetProducts(ctx context.Context, filter models.ProductFilter) ([]models.Product, error) {
	query := `
		SELECT id, name, description, price, stock_quantity, category_id, sku, image_url, is_active, created_at, updated_at
		FROM products
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, *filter.CategoryID)
		argIndex++
	}

	if filter.MinPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", argIndex)
		args = append(args, *filter.MinPrice)
		argIndex++
	}

	if filter.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", argIndex)
		args = append(args, *filter.MaxPrice)
		argIndex++
	}

	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1)
		args = append(args, "%"+filter.Search+"%", "%"+filter.Search+"%")
		argIndex += 2
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.CategoryID,
			&product.SKU,
			&product.ImageURL,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *Repository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock_quantity, category_id, sku, image_url, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.CategoryID,
		&product.SKU,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, id int, req models.UpdateProductRequest) error {
	query := `
		UPDATE products
		SET 
			name = COALESCE($2, name),
			description = COALESCE($3, description),
			price = COALESCE($4, price),
			stock_quantity = COALESCE($5, stock_quantity),
			category_id = COALESCE($6, category_id),
			image_url = COALESCE($7, image_url),
			is_active = COALESCE($8, is_active),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Name,
		req.Description,
		req.Price,
		req.StockQuantity,
		req.CategoryID,
		req.ImageURL,
		req.IsActive,
	)

	return err
}

func (r *Repository) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

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

func (r *Repository) AddToCart(ctx context.Context, userID int, req models.AddToCartRequest) (*models.CartItem, error) {
	query := `
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id) 
		DO UPDATE SET quantity = cart_items.quantity + $3, updated_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, product_id, quantity, created_at, updated_at
	`

	var cartItem models.CartItem
	err := r.db.QueryRow(ctx, query, userID, req.ProductID, req.Quantity).Scan(
		&cartItem.ID,
		&cartItem.UserID,
		&cartItem.ProductID,
		&cartItem.Quantity,
		&cartItem.CreatedAt,
		&cartItem.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	product, err := r.GetProductByID(ctx, req.ProductID)
	if err == nil {
		cartItem.Product = product
	}

	return &cartItem, nil
}

func (r *Repository) GetCartItems(ctx context.Context, userID int) ([]models.CartItem, error) {
	query := `
		SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
		       p.id, p.name, p.description, p.price, p.stock_quantity, p.category_id, p.sku, p.image_url, p.is_active, p.created_at, p.updated_at
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []models.CartItem
	for rows.Next() {
		var cartItem models.CartItem
		var product models.Product
		err := rows.Scan(
			&cartItem.ID,
			&cartItem.UserID,
			&cartItem.ProductID,
			&cartItem.Quantity,
			&cartItem.CreatedAt,
			&cartItem.UpdatedAt,
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
		cartItem.Product = &product
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}

func (r *Repository) UpdateCartItem(ctx context.Context, userID, productID int, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = $3, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND product_id = $2
	`

	_, err := r.db.Exec(ctx, query, userID, productID, quantity)
	return err
}

func (r *Repository) RemoveFromCart(ctx context.Context, userID, productID int) error {
	query := `DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`
	_, err := r.db.Exec(ctx, query, userID, productID)
	return err
}

func (r *Repository) ClearCart(ctx context.Context, userID int) error {
	query := `DELETE FROM cart_items WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *Repository) CreateOrder(ctx context.Context, req models.CreateOrderRequest, userID int) (*models.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	orderNumber := "ORD-" + uuid.New().String()[:8]

	var totalAmount float64
	for _, item := range req.Items {
		product, err := r.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found: %w", item.ProductID, err)
		}
		totalAmount += float64(item.Quantity) * product.Price
	}

	orderQuery := `
		INSERT INTO orders (user_id, order_number, status, total_amount, shipping_address, billing_address)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, order_number, status, total_amount, shipping_address, billing_address, created_at, updated_at
	`

	var order models.Order
	err = tx.QueryRow(ctx, orderQuery,
		userID,
		orderNumber,
		"pending",
		totalAmount,
		req.ShippingAddress,
		req.BillingAddress,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.BillingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		product, err := r.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found: %w", item.ProductID, err)
		}

		unitPrice := product.Price
		totalPrice := float64(item.Quantity) * unitPrice

		itemQuery := `
			INSERT INTO order_items (order_id, product_id, quantity, unit_price, total_price)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, order_id, product_id, quantity, unit_price, total_price
		`

		var orderItem models.OrderItem
		err = tx.QueryRow(ctx, itemQuery,
			order.ID,
			item.ProductID,
			item.Quantity,
			unitPrice,
			totalPrice,
		).Scan(
			&orderItem.ID,
			&orderItem.OrderID,
			&orderItem.ProductID,
			&orderItem.Quantity,
			&orderItem.UnitPrice,
			&orderItem.TotalPrice,
		)

		if err != nil {
			return nil, err
		}

		order.OrderItems = append(order.OrderItems, orderItem)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *Repository) GetOrdersByUserID(ctx context.Context, userID int) ([]models.Order, error) {
	query := `
		SELECT id, user_id, order_number, status, total_amount, shipping_address, billing_address, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.Status,
			&order.TotalAmount,
			&order.ShippingAddress,
			&order.BillingAddress,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *Repository) GetOrderByID(ctx context.Context, id int) (*models.Order, error) {
	query := `
		SELECT id, user_id, order_number, status, total_amount, shipping_address, billing_address, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var order models.Order
	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.BillingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	itemsQuery := `
		SELECT id, order_id, product_id, quantity, unit_price, total_price
		FROM order_items
		WHERE order_id = $1
	`

	itemRows, err := r.db.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item models.OrderItem
		err := itemRows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalPrice,
		)
		if err != nil {
			return nil, err
		}
		order.OrderItems = append(order.OrderItems, item)
	}

	return &order, nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, id int, status string) error {
	query := `
		UPDATE orders
		SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id, status)
	return err
}

func (r *Repository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	query := `
		SELECT id, user_id, order_number, status, total_amount, shipping_address, billing_address, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.Status,
			&order.TotalAmount,
			&order.ShippingAddress,
			&order.BillingAddress,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *Repository) CreatePayment(ctx context.Context, payment models.CreatePaymentRequest) (*models.Payment, error) {
	query := `
		INSERT INTO payments (order_id, payment_method, payment_status, amount)
		VALUES ($1, $2, $3, $4)
		RETURNING id, order_id, payment_method, payment_status, amount, transaction_id, created_at, updated_at
	`

	var paymentRecord models.Payment
	err := r.db.QueryRow(ctx, query,
		payment.OrderID,
		payment.PaymentMethod,
		"pending",
		0, // Will be updated with actual amount
	).Scan(
		&paymentRecord.ID,
		&paymentRecord.OrderID,
		&paymentRecord.PaymentMethod,
		&paymentRecord.PaymentStatus,
		&paymentRecord.Amount,
		&paymentRecord.TransactionID,
		&paymentRecord.CreatedAt,
		&paymentRecord.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &paymentRecord, nil
}

func (r *Repository) UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error {
	query := `
		UPDATE payments
		SET payment_status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, paymentID, status)
	return err
}

func (r *Repository) CreateReview(ctx context.Context, review models.CreateReviewRequest, userID int) (*models.ProductReview, error) {
	query := `
		INSERT INTO product_reviews (product_id, user_id, rating, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, product_id, user_id, rating, comment, created_at, updated_at
	`

	var productReview models.ProductReview
	err := r.db.QueryRow(ctx, query,
		review.ProductID,
		userID,
		review.Rating,
		review.Comment,
	).Scan(
		&productReview.ID,
		&productReview.ProductID,
		&productReview.UserID,
		&productReview.Rating,
		&productReview.Comment,
		&productReview.CreatedAt,
		&productReview.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &productReview, nil
}

func (r *Repository) GetProductReviews(ctx context.Context, productID int) ([]models.ProductReview, error) {
	query := `
		SELECT pr.id, pr.product_id, pr.user_id, pr.rating, pr.comment, pr.created_at, pr.updated_at,
		       u.id, u.email, u.first_name, u.last_name
		FROM product_reviews pr
		JOIN users u ON pr.user_id = u.id
		WHERE pr.product_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ProductReview
	for rows.Next() {
		var review models.ProductReview
		var user models.User
		err := rows.Scan(
			&review.ID,
			&review.ProductID,
			&review.UserID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, err
		}
		review.User = &user
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *Repository) UpdateReview(ctx context.Context, reviewID, userID int, req models.UpdateReviewRequest) (*models.ProductReview, error) {
	query := `
		UPDATE product_reviews
		SET rating = COALESCE($3, rating),
		    comment = COALESCE($4, comment),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
		RETURNING id, product_id, user_id, rating, comment, created_at, updated_at
	`

	var review models.ProductReview
	err := r.db.QueryRow(ctx, query,
		reviewID,
		userID,
		req.Rating,
		req.Comment,
	).Scan(
		&review.ID,
		&review.ProductID,
		&review.UserID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (r *Repository) DeleteReview(ctx context.Context, reviewID, userID int) error {
	query := `DELETE FROM product_reviews WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, reviewID, userID)
	return err
}

func (r *Repository) GetUserReview(ctx context.Context, productID, userID int) (*models.ProductReview, error) {
	query := `
		SELECT id, product_id, user_id, rating, comment, created_at, updated_at
		FROM product_reviews
		WHERE product_id = $1 AND user_id = $2
	`

	var review models.ProductReview
	err := r.db.QueryRow(ctx, query, productID, userID).Scan(
		&review.ID,
		&review.ProductID,
		&review.UserID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (r *Repository) GetProductAverageRating(ctx context.Context, productID int) (float64, int, error) {
	query := `
		SELECT COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as review_count
		FROM product_reviews
		WHERE product_id = $1
	`

	var avgRating float64
	var reviewCount int
	err := r.db.QueryRow(ctx, query, productID).Scan(&avgRating, &reviewCount)
	if err != nil {
		return 0, 0, err
	}

	return avgRating, reviewCount, nil
}

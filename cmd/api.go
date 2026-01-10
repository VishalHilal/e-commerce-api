package main

import (
	"log"
	"net/http"
	"time"

	"github.com/VishalHilal/e-commerce-api/internal/adapters/postgresql"
	"github.com/VishalHilal/e-commerce-api/internal/auth"
	"github.com/VishalHilal/e-commerce-api/internal/cart"
	"github.com/VishalHilal/e-commerce-api/internal/orders"
	"github.com/VishalHilal/e-commerce-api/internal/products"
	"github.com/VishalHilal/e-commerce-api/internal/reviews"
	"github.com/VishalHilal/e-commerce-api/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	repo := postgresql.New(app.db)
	jwtSvc := auth.NewJWTService("your-secret-key-change-in-production")

	userService := users.NewService(repo, jwtSvc)
	userHandler := users.NewHandler(userService)
	r.Post("/auth/register", userHandler.Register)
	r.Post("/auth/login", userHandler.Login)
	r.Group(func(r chi.Router) {
		r.Use(jwtSvc.AuthMiddleware)
		r.Get("/auth/profile", userHandler.GetProfile)
		r.Put("/auth/profile", userHandler.UpdateProfile)
	})

	productService := products.NewService(repo)
	productHandler := products.NewHandler(productService)
	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{id}", productHandler.GetProduct)

	r.Group(func(r chi.Router) {
		r.Use(jwtSvc.AuthMiddleware, auth.RequireRole("admin"))
		r.Post("/products", productHandler.CreateProduct)
		r.Put("/products/{id}", productHandler.UpdateProduct)
		r.Delete("/products/{id}", productHandler.DeleteProduct)
	})

	cartService := cart.NewService(repo)
	cartHandler := cart.NewHandler(cartService)
	r.Group(func(r chi.Router) {
		r.Use(jwtSvc.AuthMiddleware)
		r.Get("/cart", cartHandler.GetCart)
		r.Post("/cart", cartHandler.AddToCart)
		r.Put("/cart/{product_id}", cartHandler.UpdateCartItem)
		r.Delete("/cart/{product_id}", cartHandler.RemoveFromCart)
		r.Delete("/cart", cartHandler.ClearCart)
	})

	orderService := orders.NewService(repo)
	orderHandler := orders.NewHandler(orderService)
	r.Group(func(r chi.Router) {
		r.Use(jwtSvc.AuthMiddleware)
		r.Post("/orders", orderHandler.CreateOrder)
		r.Get("/orders", orderHandler.GetUserOrders)
		r.Get("/orders/{id}", orderHandler.GetOrder)
		r.Post("/payments", orderHandler.ProcessPayment)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtSvc.AuthMiddleware, auth.RequireRole("admin"))
		r.Get("/admin/orders", orderHandler.GetAllOrders)
		r.Put("/admin/orders/{id}", orderHandler.UpdateOrderStatus)
	})

	reviewService := reviews.NewService(repo)
	reviewHandler := reviews.NewHandler(reviewService)
	r.Group(func(r chi.Router) {
		r.Get("/products/{product_id}/reviews", reviewHandler.GetProductReviews)
		r.Get("/products/{product_id}/reviews-with-avg", reviewHandler.GetProductWithReviews)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtSvc.AuthMiddleware)
		r.Post("/products/{product_id}/reviews", reviewHandler.CreateReview)
		r.Put("/reviews/{id}", reviewHandler.UpdateReview)
		r.Delete("/reviews/{id}", reviewHandler.DeleteReview)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	// logger
	db *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

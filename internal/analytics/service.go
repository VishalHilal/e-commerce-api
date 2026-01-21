package analytics

import (
	"context"
	"fmt"
	"time"
)

type AnalyticsService struct {
	repo Repository
}

type Repository interface {
	GetSalesReport(ctx context.Context, startDate, endDate time.Time) (*SalesReport, error)
	GetTopProducts(ctx context.Context, limit int) ([]TopProduct, error)
	GetCustomerMetrics(ctx context.Context, startDate, endDate time.Time) (*CustomerMetrics, error)
	GetOrderStats(ctx context.Context, startDate, endDate time.Time) (*OrderStats, error)
	GetRevenueByCategory(ctx context.Context, startDate, endDate time.Time) ([]CategoryRevenue, error)
}

type SalesReport struct {
	Period            string       `json:"period"`
	TotalRevenue      float64      `json:"total_revenue"`
	TotalOrders       int          `json:"total_orders"`
	AverageOrderValue float64      `json:"average_order_value"`
	TopProducts       []TopProduct `json:"top_products"`
	DailySales        []DailySale  `json:"daily_sales"`
}

type TopProduct struct {
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	QuantitySold int     `json:"quantity_sold"`
	Revenue      float64 `json:"revenue"`
}

type DailySale struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int     `json:"orders"`
}

type CustomerMetrics struct {
	TotalCustomers     int             `json:"total_customers"`
	NewCustomers       int             `json:"new_customers"`
	ReturningCustomers int             `json:"returning_customers"`
	TopCustomers       []CustomerStats `json:"top_customers"`
}

type CustomerStats struct {
	CustomerID   int     `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	TotalOrders  int     `json:"total_orders"`
	TotalSpent   float64 `json:"total_spent"`
	AverageOrder float64 `json:"average_order"`
}

type OrderStats struct {
	TotalOrders       int     `json:"total_orders"`
	PendingOrders     int     `json:"pending_orders"`
	CompletedOrders   int     `json:"completed_orders"`
	CancelledOrders   int     `json:"cancelled_orders"`
	TotalRevenue      float64 `json:"total_revenue"`
	AverageOrderValue float64 `json:"average_order_value"`
}

type CategoryRevenue struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Revenue      float64 `json:"revenue"`
	OrderCount   int     `json:"order_count"`
}

func NewAnalyticsService(repo Repository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) GetSalesReport(ctx context.Context, startDate, endDate time.Time) (*SalesReport, error) {
	report, err := s.repo.GetSalesReport(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales report: %w", err)
	}
	return report, nil
}

func (s *AnalyticsService) GetTopProducts(ctx context.Context, limit int) ([]TopProduct, error) {
	products, err := s.repo.GetTopProducts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top products: %w", err)
	}
	return products, nil
}

func (s *AnalyticsService) GetCustomerMetrics(ctx context.Context, startDate, endDate time.Time) (*CustomerMetrics, error) {
	metrics, err := s.repo.GetCustomerMetrics(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer metrics: %w", err)
	}
	return metrics, nil
}

func (s *AnalyticsService) GetOrderStats(ctx context.Context, startDate, endDate time.Time) (*OrderStats, error) {
	stats, err := s.repo.GetOrderStats(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get order stats: %w", err)
	}
	return stats, nil
}

func (s *AnalyticsService) GetRevenueByCategory(ctx context.Context, startDate, endDate time.Time) ([]CategoryRevenue, error) {
	revenue, err := s.repo.GetRevenueByCategory(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue by category: %w", err)
	}
	return revenue, nil
}

// Dashboard data aggregation
func (s *AnalyticsService) GetDashboardData(ctx context.Context) (*DashboardData, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Get data for different time periods
	monthlyReport, _ := s.GetSalesReport(ctx, startOfMonth, now)

	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)
	endOfLastMonth := startOfMonth.AddDate(0, 0, -1)
	lastMonthReport, _ := s.GetSalesReport(ctx, startOfLastMonth, endOfLastMonth)

	topProducts, _ := s.GetTopProducts(ctx, 10)
	customerMetrics, _ := s.GetCustomerMetrics(ctx, startOfMonth, now)
	orderStats, _ := s.GetOrderStats(ctx, startOfMonth, now)

	return &DashboardData{
		CurrentMonth:    monthlyReport,
		PreviousMonth:   lastMonthReport,
		TopProducts:     topProducts,
		CustomerMetrics: customerMetrics,
		OrderStats:      orderStats,
		RevenueGrowth:   calculateGrowth(lastMonthReport.TotalRevenue, monthlyReport.TotalRevenue),
		OrderGrowth:     calculateGrowth(float64(lastMonthReport.TotalOrders), float64(monthlyReport.TotalOrders)),
		CustomerGrowth:  calculateGrowth(float64(lastMonthReport.NewCustomers), float64(customerMetrics.NewCustomers)),
	}, nil
}

type DashboardData struct {
	CurrentMonth    *SalesReport     `json:"current_month"`
	PreviousMonth   *SalesReport     `json:"previous_month"`
	TopProducts     []TopProduct     `json:"top_products"`
	CustomerMetrics *CustomerMetrics `json:"customer_metrics"`
	OrderStats      *OrderStats      `json:"order_stats"`
	RevenueGrowth   float64          `json:"revenue_growth"`
	OrderGrowth     float64          `json:"order_growth"`
	CustomerGrowth  float64          `json:"customer_growth"`
}

func calculateGrowth(previous, current float64) float64 {
	if previous == 0 {
		return 0
	}
	return ((current - previous) / previous) * 100
}

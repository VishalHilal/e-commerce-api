package health

import (
	"context"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/VishalHilal/e-commerce-api/internal/json"
	"github.com/jackc/pgx/v5"
)

type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Checks    map[string]CheckResult `json:"checks"`
	Uptime    time.Duration          `json:"uptime"`
}

type CheckResult struct {
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency,omitempty"`
}

type HealthChecker struct {
	db        *pgx.Conn
	startTime time.Time
	mu        sync.RWMutex
	checks    map[string]func(ctx context.Context) CheckResult
}

func NewHealthChecker(db *pgx.Conn) *HealthChecker {
	return &HealthChecker{
		db:        db,
		startTime: time.Now(),
		checks:    make(map[string]func(ctx context.Context) CheckResult),
	}
}

func (hc *HealthChecker) AddCheck(name string, check func(ctx context.Context) CheckResult) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks[name] = check
}

func (hc *HealthChecker) CheckDatabase(ctx context.Context) CheckResult {
	start := time.Now()

	err := hc.db.Ping(ctx)
	latency := time.Since(start)

	if err != nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: "Database connection failed",
			Latency: latency,
		}
	}

	return CheckResult{
		Status:  "healthy",
		Message: "Database connection successful",
		Latency: latency,
	}
}

func (hc *HealthChecker) CheckMemory(ctx context.Context) CheckResult {
	// Simple memory check
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Check if we're using too much memory (e.g., > 1GB)
	const maxMemoryBytes = 1024 * 1024 * 1024 // 1GB

	if m.Alloc > maxMemoryBytes {
		return CheckResult{
			Status:  "unhealthy",
			Message: "High memory usage",
		}
	}

	return CheckResult{
		Status:  "healthy",
		Message: "Memory usage normal",
	}
}

func (hc *HealthChecker) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		status := HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now(),
			Version:   "1.0.0",
			Checks:    make(map[string]CheckResult),
			Uptime:    time.Since(hc.startTime),
		}

		// Run all checks
		allHealthy := true
		for name, check := range hc.checks {
			result := check(ctx)
			status.Checks[name] = result

			if result.Status != "healthy" {
				allHealthy = false
				status.Status = "unhealthy"
			}
		}

		// Set HTTP status based on overall health
		httpStatus := http.StatusOK
		if !allHealthy {
			httpStatus = http.StatusServiceUnavailable
		}

		json.Write(w, httpStatus, status)
	})
}

func (hc *HealthChecker) ReadyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Check critical services
		dbCheck := hc.CheckDatabase(ctx)

		if dbCheck.Status == "healthy" {
			json.Write(w, http.StatusOK, map[string]string{
				"status": "ready",
			})
		} else {
			json.Write(w, http.StatusServiceUnavailable, map[string]string{
				"status": "not ready",
			})
		}
	})
}

func (hc *HealthChecker) LiveHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.Write(w, http.StatusOK, map[string]string{
			"status": "alive",
		})
	})
}

package products

import (
	"net/http"
	"strconv"

	"github.com/VishalHilal/e-commerce-api/internal/auth"
	"github.com/VishalHilal/e-commerce-api/internal/json"
	"github.com/VishalHilal/e-commerce-api/internal/models"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service *Service
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	filter := models.ProductFilter{
		CategoryID: getQueryInt(r, "category_id"),
		MinPrice:   getQueryFloat(r, "min_price"),
		MaxPrice:   getQueryFloat(r, "max_price"),
		Search:     r.URL.Query().Get("search"),
		Page:       getQueryIntDefault(r, "page", 1),
		Limit:      getQueryIntDefault(r, "limit", 20),
	}

	active := true
	filter.IsActive = &active

	products, err := h.service.ListProducts(r.Context(), filter)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]interface{}{
		"products": products,
		"count":    len(products),
	})
}

func (h *handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		json.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	json.Write(w, http.StatusOK, product)
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil || claims.Role != "admin" {
		json.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	var req models.CreateProductRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.service.CreateProduct(r.Context(), req)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, product)
}

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil || claims.Role != "admin" {
		json.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req models.UpdateProductRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateProduct(r.Context(), id, req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Product updated successfully"})
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil || claims.Role != "admin" {
		json.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

func getQueryInt(r *http.Request, key string) *int {
	if val := r.URL.Query().Get(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return &intVal
		}
	}
	return nil
}

func getQueryIntDefault(r *http.Request, key string, defaultValue int) int {
	if val := getQueryInt(r, key); val != nil {
		return *val
	}
	return defaultValue
}

func getQueryFloat(r *http.Request, key string) *float64 {
	if val := r.URL.Query().Get(key); val != "" {
		if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
			return &floatVal
		}
	}
	return nil
}

package orders

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

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreateOrderRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	order, err := h.service.CreateOrder(r.Context(), req, claims.UserID)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, order)
}

func (h *handler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	orders, err := h.service.GetUserOrders(r.Context(), claims.UserID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
	})
}

func (h *handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.service.GetOrder(r.Context(), orderID, claims.UserID)
	if err != nil {
		json.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	json.Write(w, http.StatusOK, order)
}

func (h *handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil || claims.Role != "admin" {
		json.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateOrderStatus(r.Context(), orderID, req.Status); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Order status updated successfully"})
}

func (h *handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil || claims.Role != "admin" {
		json.WriteError(w, http.StatusForbidden, "Admin access required")
		return
	}

	orders, err := h.service.GetAllOrders(r.Context())
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
	})
}

func (h *handler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreatePaymentRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	payment, err := h.service.ProcessPayment(r.Context(), req)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, payment)
}

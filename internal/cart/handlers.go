package cart

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

func (h *handler) GetCart(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	cart, err := h.service.GetCart(r.Context(), claims.UserID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, cart)
}

func (h *handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.AddToCartRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cartItem, err := h.service.AddToCart(r.Context(), claims.UserID, req)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, cartItem)
}

func (h *handler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req models.UpdateCartRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateCartItem(r.Context(), claims.UserID, productID, req.Quantity); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Cart item updated successfully"})
}

func (h *handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.RemoveFromCart(r.Context(), claims.UserID, productID); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Item removed from cart successfully"})
}

func (h *handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.service.ClearCart(r.Context(), claims.UserID); err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Cart cleared successfully"})
}

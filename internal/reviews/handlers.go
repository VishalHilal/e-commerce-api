package reviews

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

func (h *handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreateReviewRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	review, err := h.service.CreateReview(r.Context(), req, claims.UserID)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, review)
}

func (h *handler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	reviews, err := h.service.GetProductReviews(r.Context(), productID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]interface{}{
		"reviews": reviews,
		"count":   len(reviews),
	})
}

func (h *handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	reviewID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid review ID")
		return
	}

	var req models.UpdateReviewRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	review, err := h.service.UpdateReview(r.Context(), reviewID, claims.UserID, req)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusOK, review)
}

func (h *handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	reviewID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid review ID")
		return
	}

	if err := h.service.DeleteReview(r.Context(), reviewID, claims.UserID); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Review deleted successfully"})
}

func (h *handler) GetProductWithReviews(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	productWithReviews, err := h.service.GetProductWithReviews(r.Context(), productID)
	if err != nil {
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, productWithReviews)
}

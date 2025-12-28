package products

import (
	"encoding/json"
	"net/http"

	"github.com/VishalHilal/e-commerce-api/internal/products"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {

	products := []string{"Hellow", "World"}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

package users

import (
	"net/http"

	"github.com/VishalHilal/e-commerce-api/internal/auth"
	"github.com/VishalHilal/e-commerce-api/internal/json"
	"github.com/VishalHilal/e-commerce-api/internal/models"
)

type handler struct {
	service *Service
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.service.Register(r.Context(), req)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, resp)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		json.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	json.Write(w, http.StatusOK, resp)
}

func (h *handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.service.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		json.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	json.Write(w, http.StatusOK, user)
}

func (h *handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetUserFromContext(r.Context())
	if claims == nil {
		json.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.User
	if err := json.Read(r, &req); err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateProfile(r.Context(), claims.UserID, req); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

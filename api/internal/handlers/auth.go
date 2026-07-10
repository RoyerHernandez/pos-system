package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	if req.Usuario == "" || req.Password == "" {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "usuario and password are required")
		return
	}

	resp, err := h.authService.Login(req.Usuario, req.Password)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
		return
	}

	RespondJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "refresh_token is required")
		return
	}

	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid refresh token")
		return
	}

	RespondJSON(w, http.StatusOK, resp)
}


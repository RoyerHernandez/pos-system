package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/RoyerHernandez/pos-system/api/internal/service"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
)

type authService interface {
	Login(username, password string) (*service.AuthResponse, error)
	RefreshToken(token string) (*service.AuthResponse, error)
}

// AuthController handles authentication endpoints.
type AuthController struct {
	service authService
}

// NewAuthController creates a new AuthController.
func NewAuthController(svc authService) (*AuthController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &AuthController{service: svc}, nil
}

// Login authenticates a user and returns access/refresh tokens.
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	usuario, password := req.toDomain()
	if usuario == "" || password == "" {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "usuario and password are required")
		return
	}

	resp, err := c.service.Login(usuario, password)
	if err != nil {
		var unauthorized apperror.UnauthorizedError
		if errors.As(err, &unauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", unauthorized.Message)
			return
		}
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
		return
	}

	RespondJSON(w, http.StatusOK, authResponseFromService(resp))
}

// Refresh exchanges a refresh token for new access/refresh tokens.
func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "refresh_token is required")
		return
	}

	resp, err := c.service.RefreshToken(req.RefreshToken)
	if err != nil {
		var unauthorized apperror.UnauthorizedError
		if errors.As(err, &unauthorized) {
			RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", unauthorized.Message)
			return
		}
		log.Printf("controller - Refresh: %v", err)
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid refresh token")
		return
	}

	RespondJSON(w, http.StatusOK, authResponseFromService(resp))
}

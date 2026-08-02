package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/go-chi/chi/v5"
)

type cashRegisterService interface {
	OpenRegister(userID int, montoApertura float64) (*domain.CashRegister, error)
	CloseRegister(registerID, userID int) (*domain.CashRegister, error)
	GetCurrentOpen(userID int) (*domain.CashRegister, error)
	GetAll(page, perPage int) ([]domain.CashRegister, int, error)
}

// CashRegisterController handles cash register endpoints.
type CashRegisterController struct {
	service cashRegisterService
}

// NewCashRegisterController creates a new CashRegisterController.
func NewCashRegisterController(svc cashRegisterService) (*CashRegisterController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &CashRegisterController{service: svc}, nil
}

// Open opens a new cash register session.
func (c *CashRegisterController) Open(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	var req OpenCashRegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	reg, err := c.service.OpenRegister(claims.UserID, req.MontoApertura)
	if err != nil {
		var validation apperror.ValidationError
		if errors.As(err, &validation) {
			RespondError(w, http.StatusBadRequest, "BAD_REQUEST", validation.Message)
			return
		}
		var conflict apperror.ConflictError
		if errors.As(err, &conflict) {
			RespondError(w, http.StatusConflict, "CONFLICT", conflict.Message)
			return
		}
		log.Printf("controller - OpenCashRegister: %v", err)
		RespondError(w, http.StatusBadRequest, "CASHREGISTER_ERROR", err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, cashRegisterFromDomain(*reg))
}

// Close closes an open cash register session.
func (c *CashRegisterController) Close(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid register id")
		return
	}

	reg, err := c.service.CloseRegister(id, claims.UserID)
	if err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		var conflict apperror.ConflictError
		if errors.As(err, &conflict) {
			RespondError(w, http.StatusConflict, "CONFLICT", conflict.Message)
			return
		}
		log.Printf("controller - CloseCashRegister: %v", err)
		RespondError(w, http.StatusBadRequest, "CASHREGISTER_ERROR", err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, cashRegisterFromDomain(*reg))
}

// GetCurrent returns the currently open register for the authenticated user.
func (c *CashRegisterController) GetCurrent(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	reg, err := c.service.GetCurrentOpen(claims.UserID)
	if err != nil {
		log.Printf("controller - GetCurrentCashRegister: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch current register")
		return
	}
	if reg == nil {
		RespondJSON(w, http.StatusOK, map[string]interface{}{"register": nil})
		return
	}

	RespondJSON(w, http.StatusOK, cashRegisterFromDomain(*reg))
}

// GetAll returns paginated cash register history.
func (c *CashRegisterController) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page := 1
	if v := q.Get("page"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid page")
			return
		}
		page = parsed
	}

	perPage := 25
	if v := q.Get("per_page"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid per_page")
			return
		}
		perPage = parsed
	}
	if perPage > 200 {
		perPage = 200
	}

	registers, total, err := c.service.GetAll(page, perPage)
	if err != nil {
		log.Printf("controller - GetAllCashRegisters: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch cash registers")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"registers": cashRegisterListFromDomain(registers),
		"total":     total,
		"page":      page,
	})
}

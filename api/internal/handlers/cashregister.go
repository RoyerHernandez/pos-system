package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/middleware"
	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/services"
	"github.com/go-chi/chi/v5"
)

type CashRegisterHandler struct {
	cashService *services.CashRegisterService
}

func NewCashRegisterHandler(cashService *services.CashRegisterService) *CashRegisterHandler {
	return &CashRegisterHandler{cashService: cashService}
}

func (h *CashRegisterHandler) Open(w http.ResponseWriter, r *http.Request) {
	claims := middleware.UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	var req models.OpenCashRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	reg, err := h.cashService.OpenRegister(claims.UserID, req.MontoApertura)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "CASHREGISTER_ERROR", err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, reg)
}

func (h *CashRegisterHandler) Close(w http.ResponseWriter, r *http.Request) {
	claims := middleware.UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid register id")
		return
	}

	reg, err := h.cashService.CloseRegister(id, claims.UserID)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "CASHREGISTER_ERROR", err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, reg)
}

func (h *CashRegisterHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	claims := middleware.UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	reg, err := h.cashService.GetCurrentOpen(claims.UserID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch current register")
		return
	}
	if reg == nil {
		RespondJSON(w, http.StatusOK, map[string]interface{}{"register": nil})
		return
	}

	RespondJSON(w, http.StatusOK, reg)
}

func (h *CashRegisterHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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

	registers, total, err := h.cashService.GetAll(page, perPage)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch cash registers")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"registers": registers,
		"total":     total,
		"page":      page,
	})
}

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

type SaleHandler struct {
	saleService *services.SaleService
}

func NewSaleHandler(saleService *services.SaleService) *SaleHandler {
	return &SaleHandler{saleService: saleService}
}

func (h *SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	var req models.CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	resp, err := h.saleService.CreateSale(req, claims.UserID)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "SALE_ERROR", err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, resp)
}

func (h *SaleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	startDate := q.Get("start_date")
	endDate := q.Get("end_date")
	status := q.Get("status")

	userID := 0
	if v := q.Get("user_id"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid user_id")
			return
		}
		userID = parsed
	}

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

	sales, total, err := h.saleService.GetAll(startDate, endDate, status, userID, page, perPage)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch sales")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"sales": sales,
		"total": total,
		"page":  page,
	})
}

func (h *SaleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid sale id")
		return
	}

	resp, err := h.saleService.GetByID(id)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch sale")
		return
	}
	if resp == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "sale not found")
		return
	}

	RespondJSON(w, http.StatusOK, resp)
}

func (h *SaleHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	claims := middleware.UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid sale id")
		return
	}

	if err := h.saleService.CancelSale(id, claims.UserID); err != nil {
		RespondError(w, http.StatusBadRequest, "CANCEL_ERROR", err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "sale cancelled"})
}

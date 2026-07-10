package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/middleware"
	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/services"
)

type InventoryHandler struct {
	inventoryService *services.InventoryService
}

func NewInventoryHandler(inventoryService *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

func (h *InventoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	productID := 0
	if v := q.Get("product_id"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid product_id")
			return
		}
		productID = parsed
	}

	tipo := q.Get("type")
	startDate := q.Get("start_date")
	endDate := q.Get("end_date")

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

	movements, total, err := h.inventoryService.GetMovements(productID, tipo, startDate, endDate, page, perPage)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch inventory movements")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"movements": movements,
		"total":     total,
		"page":      page,
	})
}

func (h *InventoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	var req models.CreateMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	movement, err := h.inventoryService.CreateManualMovement(req, claims.UserID)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVENTORY_ERROR", err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, movement)
}

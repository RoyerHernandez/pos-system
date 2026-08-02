package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/internal/service"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
)

type inventoryService interface {
	GetMovements(productID int, tipo, startDate, endDate string, page, perPage int) ([]domain.InventoryMovement, int, error)
	CreateManualMovement(req service.CreateMovementRequest, userID int) (*domain.InventoryMovement, error)
}

// InventoryController handles inventory movement endpoints.
type InventoryController struct {
	service inventoryService
}

// NewInventoryController creates a new InventoryController.
func NewInventoryController(svc inventoryService) (*InventoryController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &InventoryController{service: svc}, nil
}

// GetAll returns paginated inventory movements with optional filters.
func (c *InventoryController) GetAll(w http.ResponseWriter, r *http.Request) {
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

	movements, total, err := c.service.GetMovements(productID, tipo, startDate, endDate, page, perPage)
	if err != nil {
		log.Printf("controller - GetAllInventory: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch inventory movements")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"movements": inventoryMovementListFromDomain(movements),
		"total":     total,
		"page":      page,
	})
}

// Create records a manual inventory movement.
func (c *InventoryController) Create(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	var req CreateMovementDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	movement, err := c.service.CreateManualMovement(service.CreateMovementRequest{
		IDProducto:    req.IDProducto,
		Tipo:          req.Tipo,
		Motivo:        req.Motivo,
		Cantidad:      req.Cantidad,
		Observaciones: req.Observaciones,
	}, claims.UserID)
	if err != nil {
		var validation apperror.ValidationError
		if errors.As(err, &validation) {
			RespondError(w, http.StatusBadRequest, "BAD_REQUEST", validation.Message)
			return
		}
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
		log.Printf("controller - CreateInventoryMovement: %v", err)
		RespondError(w, http.StatusBadRequest, "INVENTORY_ERROR", err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, inventoryMovementFromDomain(*movement))
}

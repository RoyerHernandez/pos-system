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

type tableService interface {
	GetAll() ([]domain.Table, error)
	GetByID(id int) (*domain.Table, error)
	Create(t *domain.Table) (int64, error)
	Update(t *domain.Table) error
	Delete(id int) error
}

// TableController handles table CRUD endpoints.
type TableController struct {
	service tableService
}

// NewTableController creates a new TableController.
func NewTableController(svc tableService) (*TableController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &TableController{service: svc}, nil
}

// GetAll returns all tables.
func (c *TableController) GetAll(w http.ResponseWriter, r *http.Request) {
	tables, err := c.service.GetAll()
	if err != nil {
		log.Printf("controller - GetAllTables: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch tables")
		return
	}
	RespondJSON(w, http.StatusOK, tableListFromDomain(tables))
}

// GetByID returns a single table by ID.
func (c *TableController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid table id")
		return
	}

	table, err := c.service.GetByID(id)
	if err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - GetTableByID: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch table")
		return
	}
	if table == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "table not found")
		return
	}

	RespondJSON(w, http.StatusOK, tableResponseFromDomain(*table))
}

// Create creates a new table from JSON body.
func (c *TableController) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTableDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	table := &domain.Table{
		Numero:    req.Numero,
		Nombre:    req.Nombre,
		Capacidad: req.Capacidad,
	}

	id, err := c.service.Create(table)
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
		log.Printf("controller - CreateTable: %v", err)
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to create table")
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
}

// Update updates an existing table from JSON body.
func (c *TableController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid table id")
		return
	}

	var req UpdateTableDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	table := &domain.Table{
		ID:        id,
		Numero:    req.Numero,
		Nombre:    req.Nombre,
		Capacidad: req.Capacidad,
	}

	if err := c.service.Update(table); err != nil {
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
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - UpdateTable: %v", err)
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to update table")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "table updated"})
}

// Delete removes a table by ID (only if free).
func (c *TableController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid table id")
		return
	}

	if err := c.service.Delete(id); err != nil {
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
		log.Printf("controller - DeleteTable: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete table")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "table deleted"})
}

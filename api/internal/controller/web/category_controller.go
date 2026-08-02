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

type categoryService interface {
	GetAll() ([]domain.Category, error)
	GetByID(id int) (*domain.Category, error)
	Create(nombre, descripcion string) (int64, error)
	Update(id int, nombre, descripcion string) error
	Delete(id int) error
}

// CategoryController handles category CRUD endpoints.
type CategoryController struct {
	service categoryService
}

// NewCategoryController creates a new CategoryController.
func NewCategoryController(svc categoryService) (*CategoryController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &CategoryController{service: svc}, nil
}

// GetAll returns all categories.
func (c *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := c.service.GetAll()
	if err != nil {
		log.Printf("controller - GetAllCategories: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch categories")
		return
	}
	RespondJSON(w, http.StatusOK, categoryListFromDomain(categories))
}

// GetByID returns a single category by ID.
func (c *CategoryController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid category id")
		return
	}

	category, err := c.service.GetByID(id)
	if err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - GetCategoryByID: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch category")
		return
	}
	if category == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "category not found")
		return
	}

	RespondJSON(w, http.StatusOK, categoryResponseFromDomain(*category))
}

// Create creates a new category.
func (c *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	id, err := c.service.Create(req.Nombre, req.Descripcion)
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
		log.Printf("controller - CreateCategory: %v", err)
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to create category")
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
}

// Update updates an existing category.
func (c *CategoryController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid category id")
		return
	}

	var req UpdateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	if err := c.service.Update(id, req.Nombre, req.Descripcion); err != nil {
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
		log.Printf("controller - UpdateCategory: %v", err)
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to update category")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "category updated"})
}

// Delete removes a category by ID.
func (c *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid category id")
		return
	}

	if err := c.service.Delete(id); err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - DeleteCategory: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete category")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "category deleted"})
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/services"
	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch categories")
		return
	}
	RespondJSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid category id")
		return
	}

	category, err := h.categoryService.GetByID(id)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch category")
		return
	}
	if category == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "category not found")
		return
	}

	RespondJSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	id, err := h.categoryService.Create(req.Nombre, req.Descripcion)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to create category")
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid category id")
		return
	}

	var req models.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	if err := h.categoryService.Update(id, req.Nombre, req.Descripcion); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to update category")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "category updated"})
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid category id")
		return
	}

	if err := h.categoryService.Delete(id); err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete category")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "category deleted"})
}

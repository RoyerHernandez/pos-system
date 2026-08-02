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

type productService interface {
	GetAll(categoryID int) ([]domain.Product, error)
	GetByID(id int) (*domain.Product, error)
	Create(codigo string, idCategoria int, descripcion string, stock int, precioCompra, precioVenta float64) (int64, error)
	Update(id int, idCategoria int, descripcion string, precioCompra, precioVenta float64) error
	UpdateImage(id int, imagen string) error
	Delete(id int) error
}

// ProductController handles product CRUD endpoints.
type ProductController struct {
	service productService
}

// NewProductController creates a new ProductController.
func NewProductController(svc productService) (*ProductController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &ProductController{service: svc}, nil
}

// GetAll returns all products, optionally filtered by category.
func (c *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	categoryID := 0
	if v := r.URL.Query().Get("category_id"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid category_id")
			return
		}
		categoryID = parsed
	}

	products, err := c.service.GetAll(categoryID)
	if err != nil {
		log.Printf("controller - GetAllProducts: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch products")
		return
	}
	RespondJSON(w, http.StatusOK, productListFromDomain(products))
}

// GetByID returns a single product by ID.
func (c *ProductController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid product id")
		return
	}

	product, err := c.service.GetByID(id)
	if err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - GetProductByID: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch product")
		return
	}
	if product == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "product not found")
		return
	}

	RespondJSON(w, http.StatusOK, productResponseFromDomain(*product))
}

// Create creates a new product from JSON body.
func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	productID, err := c.service.Create(req.Codigo, req.IDCategoria, req.Descripcion, req.Stock, req.PrecioCompra, req.PrecioVenta)
	if err != nil {
		var validation apperror.ValidationError
		if errors.As(err, &validation) {
			RespondError(w, http.StatusBadRequest, "BAD_REQUEST", validation.Message)
			return
		}
		log.Printf("controller - CreateProduct: %v", err)
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to create product")
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{"id": productID})
}

// Update updates an existing product from JSON body.
func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid product id")
		return
	}

	var req UpdateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	if err := c.service.Update(id, req.IDCategoria, req.Descripcion, req.PrecioCompra, req.PrecioVenta); err != nil {
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
		log.Printf("controller - UpdateProduct: %v", err)
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to update product")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "product updated"})
}

// Delete removes a product by ID.
func (c *ProductController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid product id")
		return
	}

	if err := c.service.Delete(id); err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - DeleteProduct: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete product")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "product deleted"})
}

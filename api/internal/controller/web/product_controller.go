package web

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/RoyerHernandez/pos-system/api/pkg/generic"
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

// Create creates a new product from multipart form data.
func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid multipart form")
		return
	}

	codigo := r.FormValue("codigo")
	idCategoria, err := strconv.Atoi(r.FormValue("id_categoria"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid id_categoria value")
		return
	}
	descripcion := r.FormValue("descripcion")
	precioCompra, err := strconv.ParseFloat(r.FormValue("precio_compra"), 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid precio_compra value")
		return
	}
	precioVenta, err := strconv.ParseFloat(r.FormValue("precio_venta"), 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid precio_venta value")
		return
	}
	stock, err := strconv.Atoi(r.FormValue("stock"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid stock value")
		return
	}

	productID, err := c.service.Create(codigo, idCategoria, descripcion, stock, precioCompra, precioVenta)
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

	// Handle image upload
	file, header, err := r.FormFile("imagen")
	if err == nil {
		defer file.Close()
		filename, err := generic.SaveImage(file, header, "productos", strconv.Itoa(int(productID)))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to upload product image")
			return
		}
		if err := c.service.UpdateImage(int(productID), filename); err != nil {
			log.Printf("controller - CreateProduct/UpdateImage: %v", err)
			RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update product image")
			return
		}
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{"id": productID})
}

// Update updates an existing product from multipart form data.
func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid product id")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid multipart form")
		return
	}

	idCategoria, err := strconv.Atoi(r.FormValue("id_categoria"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid id_categoria value")
		return
	}
	descripcion := r.FormValue("descripcion")
	precioCompra, err := strconv.ParseFloat(r.FormValue("precio_compra"), 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid precio_compra value")
		return
	}
	precioVenta, err := strconv.ParseFloat(r.FormValue("precio_venta"), 64)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid precio_venta value")
		return
	}

	if err := c.service.Update(id, idCategoria, descripcion, precioCompra, precioVenta); err != nil {
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

	// Handle image upload
	file, header, err := r.FormFile("imagen")
	if err == nil {
		defer file.Close()
		filename, err := generic.SaveImage(file, header, "productos", strconv.Itoa(id))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to upload product image")
			return
		}
		if err := c.service.UpdateImage(id, filename); err != nil {
			log.Printf("controller - UpdateProduct/UpdateImage: %v", err)
			RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update product image")
			return
		}
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

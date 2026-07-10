package handlers

import (
	"net/http"
	"strconv"

	"github.com/RoyerHernandez/pos-system/api/internal/services"
	"github.com/RoyerHernandez/pos-system/api/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categoryID := 0
	if v := r.URL.Query().Get("category_id"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid category_id")
			return
		}
		categoryID = parsed
	}

	products, err := h.productService.GetAll(categoryID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch products")
		return
	}
	RespondJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid product id")
		return
	}

	product, err := h.productService.GetByID(id)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch product")
		return
	}
	if product == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "product not found")
		return
	}

	RespondJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	stock, err := strconv.Atoi(r.FormValue("stock"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid stock value")
		return
	}

	productID, err := h.productService.Create(idCategoria, descripcion, stock, precioCompra, precioVenta)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to create product")
		return
	}

	// Handle image upload
	file, header, err := r.FormFile("imagen")
	if err == nil {
		defer file.Close()
		filename, err := utils.SaveImage(file, header, "productos", strconv.Itoa(int(productID)))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to upload product image")
			return
		}
		if err := h.productService.UpdateImage(int(productID), filename); err != nil {
			RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update product image")
			return
		}
	}

	RespondJSON(w, http.StatusCreated, map[string]interface{}{"id": productID})
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	if err := h.productService.Update(id, idCategoria, descripcion, precioCompra, precioVenta); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to update product")
		return
	}

	// Handle image upload
	file, header, err := r.FormFile("imagen")
	if err == nil {
		defer file.Close()
		filename, err := utils.SaveImage(file, header, "productos", strconv.Itoa(id))
		if err != nil {
			RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "failed to upload product image")
			return
		}
		if err := h.productService.UpdateImage(id, filename); err != nil {
			RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update product image")
			return
		}
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "product updated"})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid product id")
		return
	}

	if err := h.productService.Delete(id); err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete product")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "product deleted"})
}

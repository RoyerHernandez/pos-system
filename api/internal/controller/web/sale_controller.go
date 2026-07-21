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
	"github.com/go-chi/chi/v5"
)

type saleService interface {
	GetAll(startDate, endDate, status string, userID, page, perPage int) ([]domain.Sale, int, error)
	GetByID(id int) (*service.SaleResponse, error)
	CreateSale(req service.SaleRequest, userID int) (*service.SaleResponse, error)
	CancelSale(saleID int, userID int) error
}

// SaleController handles sale endpoints.
type SaleController struct {
	service saleService
}

// NewSaleController creates a new SaleController.
func NewSaleController(svc saleService) (*SaleController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &SaleController{service: svc}, nil
}

// Create processes a new sale.
func (c *SaleController) Create(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	var req CreateSaleDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid request body")
		return
	}

	resp, err := c.service.CreateSale(req.toServiceRequest(), claims.UserID)
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
		log.Printf("controller - CreateSale: %v", err)
		RespondError(w, http.StatusBadRequest, "SALE_ERROR", err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, saleResponseFromService(resp))
}

// GetAll returns paginated sales with optional filters.
func (c *SaleController) GetAll(w http.ResponseWriter, r *http.Request) {
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

	sales, total, err := c.service.GetAll(startDate, endDate, status, userID, page, perPage)
	if err != nil {
		log.Printf("controller - GetAllSales: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch sales")
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"sales": saleListFromDomain(sales),
		"total": total,
		"page":  page,
	})
}

// GetByID returns a single sale with its details.
func (c *SaleController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid sale id")
		return
	}

	resp, err := c.service.GetByID(id)
	if err != nil {
		var notFound apperror.NotFoundError
		if errors.As(err, &notFound) {
			RespondError(w, http.StatusNotFound, "NOT_FOUND", notFound.Message)
			return
		}
		log.Printf("controller - GetSaleByID: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch sale")
		return
	}
	if resp == nil {
		RespondError(w, http.StatusNotFound, "NOT_FOUND", "sale not found")
		return
	}

	RespondJSON(w, http.StatusOK, saleResponseFromService(resp))
}

// Cancel cancels an existing sale.
func (c *SaleController) Cancel(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		RespondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "no user context")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "invalid sale id")
		return
	}

	if err := c.service.CancelSale(id, claims.UserID); err != nil {
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
		log.Printf("controller - CancelSale: %v", err)
		RespondError(w, http.StatusBadRequest, "CANCEL_ERROR", err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "sale cancelled"})
}

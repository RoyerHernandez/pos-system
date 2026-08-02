package web

import (
	"errors"
	"log"
	"net/http"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
)

type dashboardService interface {
	GetKPIs() (*domain.DashboardKPIs, error)
}

// DashboardController handles dashboard endpoints.
type DashboardController struct {
	service dashboardService
}

// NewDashboardController creates a new DashboardController.
func NewDashboardController(svc dashboardService) (*DashboardController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &DashboardController{service: svc}, nil
}

// GetKPIs returns dashboard key performance indicators.
func (c *DashboardController) GetKPIs(w http.ResponseWriter, r *http.Request) {
	kpis, err := c.service.GetKPIs()
	if err != nil {
		log.Printf("controller - GetKPIs: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch dashboard KPIs")
		return
	}

	RespondJSON(w, http.StatusOK, dashboardKPIsFromDomain(kpis))
}

package handlers

import (
	"net/http"

	"github.com/RoyerHernandez/pos-system/api/internal/services"
)

type DashboardHandler struct {
	dashService *services.DashboardService
}

func NewDashboardHandler(dashService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashService: dashService}
}

func (h *DashboardHandler) GetKPIs(w http.ResponseWriter, r *http.Request) {
	kpis, err := h.dashService.GetKPIs()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch dashboard KPIs")
		return
	}

	RespondJSON(w, http.StatusOK, kpis)
}

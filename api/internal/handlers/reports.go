package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"

	"github.com/RoyerHernandez/pos-system/api/internal/services"
)

type ReportsHandler struct {
	reportsService *services.ReportsService
}

func NewReportsHandler(reportsService *services.ReportsService) *ReportsHandler {
	return &ReportsHandler{reportsService: reportsService}
}

func (h *ReportsHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	startDate := q.Get("start_date")
	endDate := q.Get("end_date")

	report, err := h.reportsService.GetSalesReport(startDate, endDate)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch sales report")
		return
	}

	if wantsCSV(r) {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=sales_report.csv")
		writer := csv.NewWriter(w)
		writer.Write([]string{"total_sales", "total_revenue", "total_efectivo", "total_tarjeta", "total_transferencia", "avg_ticket"})
		writer.Write([]string{
			fmt.Sprintf("%d", report.TotalSales),
			fmt.Sprintf("%.2f", report.TotalRevenue),
			fmt.Sprintf("%.2f", report.TotalEfectivo),
			fmt.Sprintf("%.2f", report.TotalTarjeta),
			fmt.Sprintf("%.2f", report.TotalTransfer),
			fmt.Sprintf("%.2f", report.AvgTicket),
		})
		writer.Flush()
		return
	}

	RespondJSON(w, http.StatusOK, report)
}

func (h *ReportsHandler) GetProductReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	startDate := q.Get("start_date")
	endDate := q.Get("end_date")

	reports, err := h.reportsService.GetProductReport(startDate, endDate)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch product report")
		return
	}

	if wantsCSV(r) {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=product_report.csv")
		writer := csv.NewWriter(w)
		writer.Write([]string{"id", "descripcion", "total_sold", "total_revenue"})
		for _, p := range reports {
			writer.Write([]string{
				fmt.Sprintf("%d", p.ID),
				p.Descripcion,
				fmt.Sprintf("%d", p.TotalSold),
				fmt.Sprintf("%.2f", p.TotalRevenue),
			})
		}
		writer.Flush()
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{"products": reports})
}

func (h *ReportsHandler) GetClientReport(w http.ResponseWriter, r *http.Request) {
	reports, err := h.reportsService.GetClientReport()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch client report")
		return
	}

	if wantsCSV(r) {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=client_report.csv")
		writer := csv.NewWriter(w)
		writer.Write([]string{"id", "nombre", "documento", "total_compras", "num_compras"})
		for _, c := range reports {
			doc := ""
			if c.Documento != nil {
				doc = *c.Documento
			}
			writer.Write([]string{
				fmt.Sprintf("%d", c.ID),
				c.Nombre,
				doc,
				fmt.Sprintf("%.2f", c.TotalCompras),
				fmt.Sprintf("%d", c.NumCompras),
			})
		}
		writer.Flush()
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{"clients": reports})
}

func wantsCSV(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "text/csv")
}

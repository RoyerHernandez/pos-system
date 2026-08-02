package web

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
)

type reportsService interface {
	GetSalesReport(startDate, endDate string) (*domain.SalesReport, error)
	GetProductReport(startDate, endDate string) ([]domain.ProductReport, error)
	GetClientReport() ([]domain.ClientReport, error)
}

// ReportsController handles report endpoints.
type ReportsController struct {
	service reportsService
}

// NewReportsController creates a new ReportsController.
func NewReportsController(svc reportsService) (*ReportsController, error) {
	if svc == nil {
		return nil, errors.New("service must not be nil")
	}
	return &ReportsController{service: svc}, nil
}

// GetSalesReport returns a sales summary report (JSON or CSV).
func (c *ReportsController) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	startDate := q.Get("start_date")
	endDate := q.Get("end_date")

	report, err := c.service.GetSalesReport(startDate, endDate)
	if err != nil {
		log.Printf("controller - GetSalesReport: %v", err)
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

	RespondJSON(w, http.StatusOK, salesReportFromDomain(report))
}

// GetProductReport returns a product performance report (JSON or CSV).
func (c *ReportsController) GetProductReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	startDate := q.Get("start_date")
	endDate := q.Get("end_date")

	reports, err := c.service.GetProductReport(startDate, endDate)
	if err != nil {
		log.Printf("controller - GetProductReport: %v", err)
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

	RespondJSON(w, http.StatusOK, map[string]interface{}{"products": productReportListFromDomain(reports)})
}

// GetClientReport returns a client purchasing report (JSON or CSV).
func (c *ReportsController) GetClientReport(w http.ResponseWriter, r *http.Request) {
	reports, err := c.service.GetClientReport()
	if err != nil {
		log.Printf("controller - GetClientReport: %v", err)
		RespondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch client report")
		return
	}

	if wantsCSV(r) {
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=client_report.csv")
		writer := csv.NewWriter(w)
		writer.Write([]string{"id", "nombre", "documento", "total_compras", "num_compras"})
		for _, cl := range reports {
			doc := ""
			if cl.Documento != nil {
				doc = *cl.Documento
			}
			writer.Write([]string{
				fmt.Sprintf("%d", cl.ID),
				cl.Nombre,
				doc,
				fmt.Sprintf("%.2f", cl.TotalCompras),
				fmt.Sprintf("%d", cl.NumCompras),
			})
		}
		writer.Flush()
		return
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{"clients": clientReportListFromDomain(reports)})
}

// wantsCSV checks if the client prefers CSV format via the Accept header.
func wantsCSV(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "text/csv")
}

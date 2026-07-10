package services

import (
	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
)

type ReportsService struct {
	reportsRepo *repositories.ReportsRepository
}

func NewReportsService(reportsRepo *repositories.ReportsRepository) *ReportsService {
	return &ReportsService{reportsRepo: reportsRepo}
}

func (s *ReportsService) GetSalesReport(startDate, endDate string) (*models.SalesReport, error) {
	return s.reportsRepo.GetSalesReport(startDate, endDate)
}

func (s *ReportsService) GetProductReport(startDate, endDate string) ([]models.ProductReport, error) {
	return s.reportsRepo.GetProductReport(startDate, endDate)
}

func (s *ReportsService) GetClientReport() ([]models.ClientReport, error) {
	return s.reportsRepo.GetClientReport()
}

package service

import (
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
)

type reportsRepository interface {
	GetSalesReport(startDate, endDate string) (*domain.SalesReport, error)
	GetProductReport(startDate, endDate string) ([]domain.ProductReport, error)
	GetClientReport() ([]domain.ClientReport, error)
}

type ReportsService struct {
	repo reportsRepository
}

func NewReportsService(repo reportsRepository) (*ReportsService, error) {
	if repo == nil {
		return nil, errors.New("reportsRepository must not be nil")
	}
	return &ReportsService{repo: repo}, nil
}

func (s *ReportsService) GetSalesReport(startDate, endDate string) (*domain.SalesReport, error) {
	return s.repo.GetSalesReport(startDate, endDate)
}

func (s *ReportsService) GetProductReport(startDate, endDate string) ([]domain.ProductReport, error) {
	return s.repo.GetProductReport(startDate, endDate)
}

func (s *ReportsService) GetClientReport() ([]domain.ClientReport, error) {
	return s.repo.GetClientReport()
}

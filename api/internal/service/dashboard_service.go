package service

import (
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
)

type dashboardRepository interface {
	GetTodaySalesCount() (int, error)
	GetTodayRevenue() (float64, error)
	GetTotalProducts() (int, error)
	GetLowStockCount() (int, error)
	GetTopProducts(limit int) ([]domain.TopProduct, error)
	GetRecentSales(limit int) ([]domain.RecentSale, error)
}

type DashboardService struct {
	repo dashboardRepository
}

func NewDashboardService(repo dashboardRepository) (*DashboardService, error) {
	if repo == nil {
		return nil, errors.New("dashboardRepository must not be nil")
	}
	return &DashboardService{repo: repo}, nil
}

func (s *DashboardService) GetKPIs() (*domain.DashboardKPIs, error) {
	salesCount, err := s.repo.GetTodaySalesCount()
	if err != nil {
		return nil, err
	}

	revenue, err := s.repo.GetTodayRevenue()
	if err != nil {
		return nil, err
	}

	totalProducts, err := s.repo.GetTotalProducts()
	if err != nil {
		return nil, err
	}

	lowStock, err := s.repo.GetLowStockCount()
	if err != nil {
		return nil, err
	}

	topProducts, err := s.repo.GetTopProducts(5)
	if err != nil {
		return nil, err
	}

	recentSales, err := s.repo.GetRecentSales(10)
	if err != nil {
		return nil, err
	}

	return &domain.DashboardKPIs{
		TodaySalesCount: salesCount,
		TodayRevenue:    revenue,
		TotalProducts:   totalProducts,
		LowStockCount:   lowStock,
		TopProducts:     topProducts,
		RecentSales:     recentSales,
	}, nil
}

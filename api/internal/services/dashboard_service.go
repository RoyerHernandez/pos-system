package services

import (
	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
)

type DashboardService struct {
	dashRepo *repositories.DashboardRepository
}

func NewDashboardService(dashRepo *repositories.DashboardRepository) *DashboardService {
	return &DashboardService{dashRepo: dashRepo}
}

func (s *DashboardService) GetKPIs() (*models.DashboardKPIs, error) {
	salesCount, err := s.dashRepo.GetTodaySalesCount()
	if err != nil {
		return nil, err
	}

	revenue, err := s.dashRepo.GetTodayRevenue()
	if err != nil {
		return nil, err
	}

	totalProducts, err := s.dashRepo.GetTotalProducts()
	if err != nil {
		return nil, err
	}

	lowStock, err := s.dashRepo.GetLowStockCount()
	if err != nil {
		return nil, err
	}

	topProducts, err := s.dashRepo.GetTopProducts(5)
	if err != nil {
		return nil, err
	}

	recentSales, err := s.dashRepo.GetRecentSales(10)
	if err != nil {
		return nil, err
	}

	return &models.DashboardKPIs{
		TodaySalesCount: salesCount,
		TodayRevenue:    revenue,
		TotalProducts:   totalProducts,
		LowStockCount:   lowStock,
		TopProducts:     topProducts,
		RecentSales:     recentSales,
	}, nil
}

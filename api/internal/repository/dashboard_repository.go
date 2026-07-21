package repository

import (
	"errors"
	"log"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type DashboardRepository struct {
	db *sqlx.DB
}

func NewDashboardRepository(db *sqlx.DB) (*DashboardRepository, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &DashboardRepository{db: db}, nil
}

func (r *DashboardRepository) GetTodaySalesCount() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM ventas WHERE DATE(fecha) = CURDATE() AND estado = 'completada'")
	if err != nil {
		log.Printf("repository - GetTodaySalesCount: %v", err)
		return 0, errors.New("failed to get today sales count")
	}
	return count, nil
}

func (r *DashboardRepository) GetTodayRevenue() (float64, error) {
	var revenue float64
	err := r.db.Get(&revenue, "SELECT COALESCE(SUM(total), 0) FROM ventas WHERE DATE(fecha) = CURDATE() AND estado = 'completada'")
	if err != nil {
		log.Printf("repository - GetTodayRevenue: %v", err)
		return 0, errors.New("failed to get today revenue")
	}
	return revenue, nil
}

func (r *DashboardRepository) GetTotalProducts() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM productos WHERE estado = 1")
	if err != nil {
		log.Printf("repository - GetTotalProducts: %v", err)
		return 0, errors.New("failed to get total products")
	}
	return count, nil
}

func (r *DashboardRepository) GetLowStockCount() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM productos WHERE estado = 1 AND stock <= stock_minimo")
	if err != nil {
		log.Printf("repository - GetLowStockCount: %v", err)
		return 0, errors.New("failed to get low stock count")
	}
	return count, nil
}

func (r *DashboardRepository) GetTopProducts(limit int) ([]domain.TopProduct, error) {
	var daos []topProductDAO
	err := r.db.Select(&daos,
		`SELECT p.id, p.descripcion, p.precio_venta, p.stock,
		        COALESCE(SUM(d.cantidad), 0) AS total_sold
		 FROM productos p
		 LEFT JOIN detalle_ventas d ON p.id = d.id_producto
		 LEFT JOIN ventas v ON d.id_venta = v.id AND v.estado = 'completada'
		 WHERE p.estado = 1
		 GROUP BY p.id
		 ORDER BY total_sold DESC
		 LIMIT ?`, limit)
	if err != nil {
		log.Printf("repository - GetTopProducts: %v", err)
		return nil, errors.New("failed to get top products")
	}
	products := make([]domain.TopProduct, len(daos))
	for i := range daos {
		products[i] = daos[i].toDomain()
	}
	return products, nil
}

func (r *DashboardRepository) GetRecentSales(limit int) ([]domain.RecentSale, error) {
	var daos []recentSaleDAO
	err := r.db.Select(&daos,
		`SELECT id, codigo_venta, total, metodo_pago, estado, fecha
		 FROM ventas
		 ORDER BY id DESC
		 LIMIT ?`, limit)
	if err != nil {
		log.Printf("repository - GetRecentSales: %v", err)
		return nil, errors.New("failed to get recent sales")
	}
	sales := make([]domain.RecentSale, len(daos))
	for i := range daos {
		sales[i] = daos[i].toDomain()
	}
	return sales, nil
}

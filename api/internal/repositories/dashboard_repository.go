package repositories

import (
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/jmoiron/sqlx"
)

type DashboardRepository struct {
	db *sqlx.DB
}

func NewDashboardRepository(db *sqlx.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetTodaySalesCount() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM ventas WHERE DATE(fecha) = CURDATE() AND estado = 'completada'")
	if err != nil {
		return 0, fmt.Errorf("today sales count: %w", err)
	}
	return count, nil
}

func (r *DashboardRepository) GetTodayRevenue() (float64, error) {
	var revenue float64
	err := r.db.Get(&revenue, "SELECT COALESCE(SUM(total), 0) FROM ventas WHERE DATE(fecha) = CURDATE() AND estado = 'completada'")
	if err != nil {
		return 0, fmt.Errorf("today revenue: %w", err)
	}
	return revenue, nil
}

func (r *DashboardRepository) GetTotalProducts() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM productos WHERE estado = 1")
	if err != nil {
		return 0, fmt.Errorf("total products: %w", err)
	}
	return count, nil
}

func (r *DashboardRepository) GetLowStockCount() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM productos WHERE estado = 1 AND stock <= stock_minimo")
	if err != nil {
		return 0, fmt.Errorf("low stock count: %w", err)
	}
	return count, nil
}

func (r *DashboardRepository) GetTopProducts(limit int) ([]models.TopProduct, error) {
	var products []models.TopProduct
	err := r.db.Select(&products,
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
		return nil, fmt.Errorf("top products: %w", err)
	}
	return products, nil
}

func (r *DashboardRepository) GetRecentSales(limit int) ([]models.RecentSale, error) {
	var sales []models.RecentSale
	err := r.db.Select(&sales,
		`SELECT id, codigo_venta, total, metodo_pago, estado, fecha
		 FROM ventas
		 ORDER BY id DESC
		 LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("recent sales: %w", err)
	}
	return sales, nil
}

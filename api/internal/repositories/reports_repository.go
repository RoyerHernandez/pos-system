package repositories

import (
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/jmoiron/sqlx"
)

type ReportsRepository struct {
	db *sqlx.DB
}

func NewReportsRepository(db *sqlx.DB) *ReportsRepository {
	return &ReportsRepository{db: db}
}

func (r *ReportsRepository) GetSalesReport(startDate, endDate string) (*models.SalesReport, error) {
	query := `SELECT
		COUNT(*) AS total_sales,
		COALESCE(SUM(total), 0) AS total_revenue,
		COALESCE(SUM(CASE WHEN metodo_pago = 'Efectivo' THEN total ELSE 0 END), 0) AS total_efectivo,
		COALESCE(SUM(CASE WHEN metodo_pago = 'Tarjeta' THEN total ELSE 0 END), 0) AS total_tarjeta,
		COALESCE(SUM(CASE WHEN metodo_pago = 'Transferencia' THEN total ELSE 0 END), 0) AS total_transferencia,
		COALESCE(AVG(total), 0) AS avg_ticket
	FROM ventas
	WHERE estado = 'completada'`

	var args []interface{}
	if startDate != "" {
		query += " AND fecha >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		query += " AND fecha <= ?"
		args = append(args, endDate+" 23:59:59")
	}

	var report models.SalesReport
	if err := r.db.Get(&report, query, args...); err != nil {
		return nil, fmt.Errorf("sales report: %w", err)
	}
	return &report, nil
}

func (r *ReportsRepository) GetProductReport(startDate, endDate string) ([]models.ProductReport, error) {
	subWhere := "WHERE v.estado = 'completada'"
	var args []interface{}

	if startDate != "" {
		subWhere += " AND v.fecha >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		subWhere += " AND v.fecha <= ?"
		args = append(args, endDate+" 23:59:59")
	}

	query := fmt.Sprintf(`SELECT p.id, p.descripcion,
		COALESCE(SUM(d.cantidad), 0) AS total_sold,
		COALESCE(SUM(d.subtotal), 0) AS total_revenue
	FROM productos p
	LEFT JOIN (
		SELECT d.* FROM detalle_ventas d
		INNER JOIN ventas v ON d.id_venta = v.id
		%s
	) d ON p.id = d.id_producto
	GROUP BY p.id ORDER BY total_revenue DESC`, subWhere)

	var reports []models.ProductReport
	if err := r.db.Select(&reports, query, args...); err != nil {
		return nil, fmt.Errorf("product report: %w", err)
	}
	return reports, nil
}

func (r *ReportsRepository) GetClientReport() ([]models.ClientReport, error) {
	var reports []models.ClientReport
	err := r.db.Select(&reports,
		`SELECT c.id, c.nombre, c.documento, c.total_compras,
		        COUNT(v.id) AS num_compras
		 FROM clientes c
		 LEFT JOIN ventas v ON c.id = v.id_cliente AND v.estado = 'completada'
		 WHERE c.estado = 1
		 GROUP BY c.id
		 ORDER BY c.total_compras DESC`)
	if err != nil {
		return nil, fmt.Errorf("client report: %w", err)
	}
	return reports, nil
}

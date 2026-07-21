package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ReportsRepository struct {
	db *sqlx.DB
}

func NewReportsRepository(db *sqlx.DB) (*ReportsRepository, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &ReportsRepository{db: db}, nil
}

func (r *ReportsRepository) GetSalesReport(startDate, endDate string) (*domain.SalesReport, error) {
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

	var dao saleReportDAO
	if err := r.db.Get(&dao, query, args...); err != nil {
		log.Printf("repository - GetSalesReport: %v", err)
		return nil, errors.New("failed to get sales report")
	}
	report := dao.toDomain()
	return &report, nil
}

func (r *ReportsRepository) GetProductReport(startDate, endDate string) ([]domain.ProductReport, error) {
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

	var daos []productReportDAO
	if err := r.db.Select(&daos, query, args...); err != nil {
		log.Printf("repository - GetProductReport: %v", err)
		return nil, errors.New("failed to get product report")
	}
	reports := make([]domain.ProductReport, len(daos))
	for i := range daos {
		reports[i] = daos[i].toDomain()
	}
	return reports, nil
}

func (r *ReportsRepository) GetClientReport() ([]domain.ClientReport, error) {
	var daos []clientReportDAO
	err := r.db.Select(&daos,
		`SELECT c.id, c.nombre, c.documento, c.total_compras,
		        COUNT(v.id) AS num_compras
		 FROM clientes c
		 LEFT JOIN ventas v ON c.id = v.id_cliente AND v.estado = 'completada'
		 WHERE c.estado = 1
		 GROUP BY c.id
		 ORDER BY c.total_compras DESC`)
	if err != nil {
		log.Printf("repository - GetClientReport: %v", err)
		return nil, errors.New("failed to get client report")
	}
	reports := make([]domain.ClientReport, len(daos))
	for i := range daos {
		reports[i] = daos[i].toDomain()
	}
	return reports, nil
}

package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/jmoiron/sqlx"
)

type SaleRepository struct {
	db *sqlx.DB
}

func NewSaleRepository(db *sqlx.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) FindAll(startDate, endDate, status string, userID, page, perPage int) ([]models.Sale, int, error) {
	query := "SELECT * FROM ventas WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM ventas WHERE 1=1"
	var args []interface{}

	if startDate != "" {
		query += " AND fecha >= ?"
		countQuery += " AND fecha >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		query += " AND fecha <= ?"
		countQuery += " AND fecha <= ?"
		args = append(args, endDate+" 23:59:59")
	}
	if status != "" {
		query += " AND estado = ?"
		countQuery += " AND estado = ?"
		args = append(args, status)
	}
	if userID > 0 {
		query += " AND id_usuario = ?"
		countQuery += " AND id_usuario = ?"
		args = append(args, userID)
	}

	var total int
	if err := r.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 25
	}
	offset := (page - 1) * perPage
	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, perPage, offset)

	var sales []models.Sale
	if err := r.db.Select(&sales, query, args...); err != nil {
		return nil, 0, err
	}
	return sales, total, nil
}

func (r *SaleRepository) FindByID(id int) (*models.Sale, error) {
	var sale models.Sale
	err := r.db.Get(&sale, "SELECT * FROM ventas WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepository) FindDetailsBySaleID(saleID int) ([]models.SaleDetail, error) {
	var details []models.SaleDetail
	err := r.db.Select(&details, "SELECT * FROM detalle_ventas WHERE id_venta = ?", saleID)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (r *SaleRepository) GetLastSaleCode(tx *sqlx.Tx) (string, error) {
	// Use MAX + CAST to lock the full index and extract the max sale number directly.
	// This serializes concurrent inserts and avoids duplicate code generation.
	var maxNum int
	err := tx.Get(&maxNum, "SELECT COALESCE(MAX(CAST(SUBSTRING(codigo_venta, 5) AS UNSIGNED)), 0) FROM ventas FOR UPDATE")
	if err != nil {
		return "", err
	}
	if maxNum == 0 {
		return "", nil
	}
	return fmt.Sprintf("VTA-%04d", maxNum), nil
}

func (r *SaleRepository) InsertSale(tx *sqlx.Tx, sale *models.Sale) (int64, error) {
	result, err := tx.Exec(
		`INSERT INTO ventas (id_usuario, id_cliente, id_caja, codigo_venta, subtotal, impuesto, descuento, total, metodo_pago, estado)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sale.IDUsuario, sale.IDCliente, sale.IDCaja, sale.CodigoVenta,
		sale.Subtotal, sale.Impuesto, sale.Descuento, sale.Total,
		sale.MetodoPago, sale.Estado,
	)
	if err != nil {
		return 0, fmt.Errorf("insert sale: %w", err)
	}
	return result.LastInsertId()
}

func (r *SaleRepository) InsertSaleDetail(tx *sqlx.Tx, detail *models.SaleDetail) error {
	_, err := tx.Exec(
		`INSERT INTO detalle_ventas (id_venta, id_producto, cantidad, precio_unitario, descuento, subtotal)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		detail.IDVenta, detail.IDProducto, detail.Cantidad,
		detail.PrecioUnitario, detail.Descuento, detail.Subtotal,
	)
	if err != nil {
		return fmt.Errorf("insert sale detail: %w", err)
	}
	return nil
}

func (r *SaleRepository) UpdateSaleStatus(tx *sqlx.Tx, id int, status string) error {
	_, err := tx.Exec("UPDATE ventas SET estado = ? WHERE id = ?", status, id)
	if err != nil {
		return fmt.Errorf("update sale status: %w", err)
	}
	return nil
}

func (r *SaleRepository) FindByIDForUpdate(tx *sqlx.Tx, id int) (*models.Sale, error) {
	var sale models.Sale
	err := tx.Get(&sale, "SELECT * FROM ventas WHERE id = ? FOR UPDATE", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *SaleRepository) FindDetailsBySaleIDTx(tx *sqlx.Tx, saleID int) ([]models.SaleDetail, error) {
	var details []models.SaleDetail
	err := tx.Select(&details, "SELECT * FROM detalle_ventas WHERE id_venta = ?", saleID)
	if err != nil {
		return nil, err
	}
	return details, nil
}

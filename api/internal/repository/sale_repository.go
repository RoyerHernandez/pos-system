package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

type SaleRepository struct {
	db *sqlx.DB
}

func NewSaleRepository(db *sqlx.DB) (*SaleRepository, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &SaleRepository{db: db}, nil
}

func (r *SaleRepository) FindAll(startDate, endDate, status string, userID, page, perPage int) ([]domain.Sale, int, error) {
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
		log.Printf("repository - FindAll(sales) count: %v", err)
		return nil, 0, errors.New("failed to count sales")
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

	var daos []saleDAO
	if err := r.db.Select(&daos, query, args...); err != nil {
		log.Printf("repository - FindAll(sales): %v", err)
		return nil, 0, errors.New("failed to find sales")
	}
	sales := make([]domain.Sale, len(daos))
	for i := range daos {
		sales[i] = daos[i].toDomain()
	}
	return sales, total, nil
}

func (r *SaleRepository) FindByID(id int) (*domain.Sale, error) {
	var dao saleDAO
	err := r.db.Get(&dao, "SELECT * FROM ventas WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFoundError{Message: "sale not found"}
	}
	if err != nil {
		log.Printf("repository - FindByID(sale): %v", err)
		return nil, errors.New("failed to find sale by id")
	}
	s := dao.toDomain()
	return &s, nil
}

func (r *SaleRepository) FindDetailsBySaleID(saleID int) ([]domain.SaleDetail, error) {
	var daos []saleDetailDAO
	err := r.db.Select(&daos, "SELECT * FROM detalle_ventas WHERE id_venta = ?", saleID)
	if err != nil {
		log.Printf("repository - FindDetailsBySaleID: %v", err)
		return nil, errors.New("failed to find sale details")
	}
	details := make([]domain.SaleDetail, len(daos))
	for i := range daos {
		details[i] = daos[i].toDomain()
	}
	return details, nil
}

func (r *SaleRepository) GetLastSaleCode(tx *sqlx.Tx) (string, error) {
	var maxNum int
	err := tx.Get(&maxNum, "SELECT COALESCE(MAX(CAST(SUBSTRING(codigo_venta, 5) AS UNSIGNED)), 0) FROM ventas FOR UPDATE")
	if err != nil {
		log.Printf("repository - GetLastSaleCode: %v", err)
		return "", errors.New("failed to get last sale code")
	}
	if maxNum == 0 {
		return "", nil
	}
	return fmt.Sprintf("VTA-%04d", maxNum), nil
}

func (r *SaleRepository) InsertSale(tx *sqlx.Tx, sale *domain.Sale) (int64, error) {
	dao := toSaleDAO(*sale)
	result, err := tx.Exec(
		`INSERT INTO ventas (id_usuario, id_cliente, id_caja, codigo_venta, subtotal, impuesto, descuento, total, metodo_pago, estado)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		dao.IDUsuario, dao.IDCliente, dao.IDCaja, dao.CodigoVenta,
		dao.Subtotal, dao.Impuesto, dao.Descuento, dao.Total,
		dao.MetodoPago, dao.Estado,
	)
	if err != nil {
		log.Printf("repository - InsertSale: %v", err)
		return 0, errors.New("failed to insert sale")
	}
	return result.LastInsertId()
}

func (r *SaleRepository) InsertSaleDetail(tx *sqlx.Tx, detail *domain.SaleDetail) error {
	dao := toSaleDetailDAO(*detail)
	_, err := tx.Exec(
		`INSERT INTO detalle_ventas (id_venta, id_producto, cantidad, precio_unitario, descuento, subtotal)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		dao.IDVenta, dao.IDProducto, dao.Cantidad,
		dao.PrecioUnitario, dao.Descuento, dao.Subtotal,
	)
	if err != nil {
		log.Printf("repository - InsertSaleDetail: %v", err)
		return errors.New("failed to insert sale detail")
	}
	return nil
}

func (r *SaleRepository) UpdateSaleStatus(tx *sqlx.Tx, id int, status string) error {
	_, err := tx.Exec("UPDATE ventas SET estado = ? WHERE id = ?", status, id)
	if err != nil {
		log.Printf("repository - UpdateSaleStatus: %v", err)
		return errors.New("failed to update sale status")
	}
	return nil
}

func (r *SaleRepository) FindByIDForUpdate(tx *sqlx.Tx, id int) (*domain.Sale, error) {
	var dao saleDAO
	err := tx.Get(&dao, "SELECT * FROM ventas WHERE id = ? FOR UPDATE", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFoundError{Message: "sale not found"}
	}
	if err != nil {
		log.Printf("repository - FindByIDForUpdate(sale): %v", err)
		return nil, errors.New("failed to find sale for update")
	}
	s := dao.toDomain()
	return &s, nil
}

func (r *SaleRepository) FindDetailsBySaleIDTx(tx *sqlx.Tx, saleID int) ([]domain.SaleDetail, error) {
	var daos []saleDetailDAO
	err := tx.Select(&daos, "SELECT * FROM detalle_ventas WHERE id_venta = ?", saleID)
	if err != nil {
		log.Printf("repository - FindDetailsBySaleIDTx: %v", err)
		return nil, errors.New("failed to find sale details in transaction")
	}
	details := make([]domain.SaleDetail, len(daos))
	for i := range daos {
		details[i] = daos[i].toDomain()
	}
	return details, nil
}

package service

import (
	"errors"
	"fmt"
	"math"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

// --- Repository interfaces used by TableOperationService ---

type tableOpTableRepo interface {
	FindByIDForUpdate(tx *sqlx.Tx, id int) (*domain.Table, error)
	UpdateState(tx *sqlx.Tx, id int, estado string, idVenta *int, idMesero *int) error
}

type tableOpSaleRepo interface {
	GetLastSaleCode(tx *sqlx.Tx) (string, error)
	InsertSale(tx *sqlx.Tx, sale *domain.Sale) (int64, error)
	InsertSaleDetail(tx *sqlx.Tx, detail *domain.SaleDetail) error
	FindByIDForUpdate(tx *sqlx.Tx, id int) (*domain.Sale, error)
	FindDetailsBySaleIDTx(tx *sqlx.Tx, saleID int) ([]domain.SaleDetail, error)
	FindDetailByIDTx(tx *sqlx.Tx, detailID int) (*domain.SaleDetail, error)
	DeleteSaleDetailTx(tx *sqlx.Tx, detailID int) error
	UpdateSaleTotals(tx *sqlx.Tx, saleID int, subtotal, total float64) error
	UpdateSaleCompleted(tx *sqlx.Tx, saleID int, metodoPago string, subtotal, total float64) error
}

type tableOpInventoryRepo interface {
	InsertMovement(tx *sqlx.Tx, movement *domain.InventoryMovement) error
}

type tableOpCashRepo interface {
	FindOpenByUserIDTx(tx *sqlx.Tx, userID int) (*domain.CashRegister, error)
	UpdateTotalsTx(tx *sqlx.Tx, cajaID int, amount float64, metodoPago string) error
}

// --- DTOs ---

// TableItemRequest represents a single product line for AddItems.
type TableItemRequest struct {
	IDProducto int
	Cantidad   int
}

// TableOrderResponse is returned by Open, AddItems, and RemoveItem.
type TableOrderResponse struct {
	Table   domain.Table
	Sale    domain.Sale
	Details []domain.SaleDetail
}

// --- Service ---

// TableOperationService handles table lifecycle: open, add items, remove item, close.
type TableOperationService struct {
	db        *sqlx.DB
	tables    tableOpTableRepo
	sales     tableOpSaleRepo
	inventory tableOpInventoryRepo
	cash      tableOpCashRepo
}

// NewTableOperationService creates a new TableOperationService.
func NewTableOperationService(
	db *sqlx.DB,
	tables tableOpTableRepo,
	sales tableOpSaleRepo,
	inventory tableOpInventoryRepo,
	cash tableOpCashRepo,
) (*TableOperationService, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	if tables == nil {
		return nil, errors.New("tableRepository must not be nil")
	}
	if sales == nil {
		return nil, errors.New("saleRepository must not be nil")
	}
	if inventory == nil {
		return nil, errors.New("inventoryRepository must not be nil")
	}
	if cash == nil {
		return nil, errors.New("cashRegisterRepository must not be nil")
	}
	return &TableOperationService{db: db, tables: tables, sales: sales, inventory: inventory, cash: cash}, nil
}

// OpenTable marks a table as ocupada and creates an empty sale with estado=abierta.
// The requesting user must have an open cash register.
func (s *TableOperationService) OpenTable(tableID, userID int) (*TableOrderResponse, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errors.New("failed to begin transaction")
	}
	defer tx.Rollback()

	// 1. Lock table row
	table, err := s.tables.FindByIDForUpdate(tx, tableID)
	if err != nil {
		return nil, err
	}
	if table.Estado != "libre" {
		return nil, apperror.ValidationError{Message: "table is not available (not libre)"}
	}

	// 2. Require open cash register
	caja, err := s.cash.FindOpenByUserIDTx(tx, userID)
	if err != nil {
		return nil, errors.New("failed to check cash register")
	}
	if caja == nil {
		return nil, apperror.ValidationError{Message: "no open cash register found for user"}
	}

	// 3. Generate sale code
	lastCode, err := s.sales.GetLastSaleCode(tx)
	if err != nil {
		return nil, errors.New("failed to generate sale code")
	}
	codigoVenta := nextSaleCode(lastCode)

	// 4. Insert sale with estado=abierta
	cajaID := caja.ID
	mesa := tableID
	sale := &domain.Sale{
		IDUsuario:   userID,
		IDCaja:      &cajaID,
		IDMesa:      &mesa,
		CodigoVenta: codigoVenta,
		MetodoPago:  "Efectivo",
		Estado:      "abierta",
	}
	saleID, err := s.sales.InsertSale(tx, sale)
	if err != nil {
		return nil, errors.New("failed to create sale")
	}
	sale.ID = int(saleID)

	// 5. Update table state
	sID := int(saleID)
	if err := s.tables.UpdateState(tx, tableID, "ocupada", &sID, &userID); err != nil {
		return nil, errors.New("failed to update table state")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	table.Estado = "ocupada"
	table.IDVentaActiva = &sID
	table.IDMesero = &userID

	return &TableOrderResponse{
		Table:   *table,
		Sale:    *sale,
		Details: nil,
	}, nil
}

// AddItems adds products to an open table's sale.
// Each product reduces stock and logs an inventory movement.
func (s *TableOperationService) AddItems(tableID, userID int, items []TableItemRequest) (*TableOrderResponse, error) {
	if len(items) == 0 {
		return nil, apperror.ValidationError{Message: "items list must not be empty"}
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errors.New("failed to begin transaction")
	}
	defer tx.Rollback()

	// 1. Lock table
	table, err := s.tables.FindByIDForUpdate(tx, tableID)
	if err != nil {
		return nil, err
	}
	if table.Estado != "ocupada" {
		return nil, apperror.ValidationError{Message: "table is not open (not ocupada)"}
	}
	if table.IDVentaActiva == nil {
		return nil, apperror.ValidationError{Message: "table has no active sale"}
	}
	saleID := *table.IDVentaActiva

	// 2. Lock sale
	sale, err := s.sales.FindByIDForUpdate(tx, saleID)
	if err != nil {
		return nil, err
	}
	if sale.Estado != "abierta" {
		return nil, apperror.ValidationError{Message: "sale is not open"}
	}

	// 3. Lock products and validate stock
	type productInfo struct {
		ID          int     `db:"id"`
		PrecioVenta float64 `db:"precio_venta"`
		Stock       int     `db:"stock"`
		Descripcion string  `db:"descripcion"`
	}
	products := make(map[int]*productInfo, len(items))
	for _, item := range items {
		if item.Cantidad <= 0 {
			return nil, apperror.ValidationError{Message: fmt.Sprintf("cantidad must be > 0 for product %d", item.IDProducto)}
		}
		var p productInfo
		if err := tx.Get(&p, "SELECT id, precio_venta, stock, descripcion FROM productos WHERE id = ? FOR UPDATE", item.IDProducto); err != nil {
			return nil, apperror.NotFoundError{Message: fmt.Sprintf("product %d not found", item.IDProducto)}
		}
		if p.Stock < item.Cantidad {
			return nil, apperror.ValidationError{Message: fmt.Sprintf("insufficient stock for product %d: available %d, requested %d", item.IDProducto, p.Stock, item.Cantidad)}
		}
		products[item.IDProducto] = &p
	}

	// 4. Insert details, update stock, log movements
	for _, item := range items {
		p := products[item.IDProducto]
		lineSubtotal := math.Round(float64(item.Cantidad)*p.PrecioVenta*100) / 100

		detail := &domain.SaleDetail{
			IDVenta:        saleID,
			IDProducto:     item.IDProducto,
			Cantidad:       item.Cantidad,
			PrecioUnitario: p.PrecioVenta,
			Subtotal:       lineSubtotal,
		}
		if err := s.sales.InsertSaleDetail(tx, detail); err != nil {
			return nil, err
		}

		if _, err := tx.Exec("UPDATE productos SET stock = stock - ? WHERE id = ?", item.Cantidad, item.IDProducto); err != nil {
			return nil, errors.New("failed to update product stock")
		}

		obs := fmt.Sprintf("Mesa %d - venta %s", tableID, sale.CodigoVenta)
		refID := saleID
		movement := &domain.InventoryMovement{
			IDProducto:    item.IDProducto,
			IDUsuario:     userID,
			Tipo:          "salida",
			Motivo:        "Venta mesa",
			Cantidad:      item.Cantidad,
			Observaciones: &obs,
			IDReferencia:  &refID,
		}
		if err := s.inventory.InsertMovement(tx, movement); err != nil {
			return nil, err
		}
	}

	// 5. Recalculate and persist sale totals
	details, err := s.sales.FindDetailsBySaleIDTx(tx, saleID)
	if err != nil {
		return nil, err
	}
	newSubtotal := calcSubtotal(details)
	newTotal := math.Round(newSubtotal*100) / 100
	if err := s.sales.UpdateSaleTotals(tx, saleID, newTotal, newTotal); err != nil {
		return nil, err
	}
	sale.Subtotal = newTotal
	sale.Total = newTotal

	if err := tx.Commit(); err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	return &TableOrderResponse{
		Table:   *table,
		Sale:    *sale,
		Details: details,
	}, nil
}

// RemoveItem removes a single detail line from an open table's sale,
// restoring stock and reversing the inventory movement.
func (s *TableOperationService) RemoveItem(tableID, detailID, userID int) (*TableOrderResponse, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errors.New("failed to begin transaction")
	}
	defer tx.Rollback()

	// 1. Lock table
	table, err := s.tables.FindByIDForUpdate(tx, tableID)
	if err != nil {
		return nil, err
	}
	if table.Estado != "ocupada" {
		return nil, apperror.ValidationError{Message: "table is not open (not ocupada)"}
	}
	if table.IDVentaActiva == nil {
		return nil, apperror.ValidationError{Message: "table has no active sale"}
	}
	saleID := *table.IDVentaActiva

	// 2. Lock sale
	sale, err := s.sales.FindByIDForUpdate(tx, saleID)
	if err != nil {
		return nil, err
	}
	if sale.Estado != "abierta" {
		return nil, apperror.ValidationError{Message: "sale is not open"}
	}

	// 3. Lock the detail row
	detail, err := s.sales.FindDetailByIDTx(tx, detailID)
	if err != nil {
		return nil, err
	}
	if detail.IDVenta != saleID {
		return nil, apperror.ValidationError{Message: "detail does not belong to this table's sale"}
	}

	// 4. Restore stock
	if _, err := tx.Exec("UPDATE productos SET stock = stock + ? WHERE id = ?", detail.Cantidad, detail.IDProducto); err != nil {
		return nil, errors.New("failed to restore product stock")
	}

	// 5. Log reversal movement
	obs := fmt.Sprintf("Eliminacion item - mesa %d - venta %s", tableID, sale.CodigoVenta)
	refID := saleID
	movement := &domain.InventoryMovement{
		IDProducto:    detail.IDProducto,
		IDUsuario:     userID,
		Tipo:          "entrada",
		Motivo:        "Eliminacion item mesa",
		Cantidad:      detail.Cantidad,
		Observaciones: &obs,
		IDReferencia:  &refID,
	}
	if err := s.inventory.InsertMovement(tx, movement); err != nil {
		return nil, err
	}

	// 6. Delete the detail
	if err := s.sales.DeleteSaleDetailTx(tx, detailID); err != nil {
		return nil, err
	}

	// 7. Recalculate totals
	remaining, err := s.sales.FindDetailsBySaleIDTx(tx, saleID)
	if err != nil {
		return nil, err
	}
	newSubtotal := calcSubtotal(remaining)
	newTotal := math.Round(newSubtotal*100) / 100
	if err := s.sales.UpdateSaleTotals(tx, saleID, newTotal, newTotal); err != nil {
		return nil, err
	}
	sale.Subtotal = newTotal
	sale.Total = newTotal

	if err := tx.Commit(); err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	return &TableOrderResponse{
		Table:   *table,
		Sale:    *sale,
		Details: remaining,
	}, nil
}

// CloseTable completes the sale, updates the cash register, and sets the table back to libre.
func (s *TableOperationService) CloseTable(tableID, userID int, metodoPago string) (*TableOrderResponse, error) {
	if metodoPago != "Efectivo" && metodoPago != "Tarjeta" && metodoPago != "Transferencia" {
		return nil, apperror.ValidationError{Message: "metodo_pago must be Efectivo, Tarjeta, or Transferencia"}
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errors.New("failed to begin transaction")
	}
	defer tx.Rollback()

	// 1. Lock table
	table, err := s.tables.FindByIDForUpdate(tx, tableID)
	if err != nil {
		return nil, err
	}
	if table.Estado != "ocupada" {
		return nil, apperror.ValidationError{Message: "table is not open (not ocupada)"}
	}
	if table.IDVentaActiva == nil {
		return nil, apperror.ValidationError{Message: "table has no active sale"}
	}
	saleID := *table.IDVentaActiva

	// 2. Lock sale
	sale, err := s.sales.FindByIDForUpdate(tx, saleID)
	if err != nil {
		return nil, err
	}
	if sale.Estado != "abierta" {
		return nil, apperror.ValidationError{Message: "sale is not open"}
	}

	// 3. Get details and final totals
	details, err := s.sales.FindDetailsBySaleIDTx(tx, saleID)
	if err != nil {
		return nil, err
	}
	subtotal := math.Round(calcSubtotal(details)*100) / 100
	total := subtotal

	// 4. Mark sale as completada
	if err := s.sales.UpdateSaleCompleted(tx, saleID, metodoPago, subtotal, total); err != nil {
		return nil, err
	}

	// 5. Update cash register
	if sale.IDCaja != nil {
		if err := s.cash.UpdateTotalsTx(tx, *sale.IDCaja, total, metodoPago); err != nil {
			return nil, err
		}
	}

	// 6. Reset table to libre
	if err := s.tables.UpdateState(tx, tableID, "libre", nil, nil); err != nil {
		return nil, errors.New("failed to reset table state")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	sale.Estado = "completada"
	sale.MetodoPago = metodoPago
	sale.Subtotal = subtotal
	sale.Total = total
	table.Estado = "libre"
	table.IDVentaActiva = nil
	table.IDMesero = nil

	return &TableOrderResponse{
		Table:   *table,
		Sale:    *sale,
		Details: details,
	}, nil
}

// --- helpers ---

func calcSubtotal(details []domain.SaleDetail) float64 {
	var sum float64
	for _, d := range details {
		sum += d.Subtotal
	}
	return sum
}

func nextSaleCode(lastCode string) string {
	if lastCode == "" {
		return "VTA-000001"
	}
	var n int
	fmt.Sscanf(lastCode, "VTA-%d", &n)
	return fmt.Sprintf("VTA-%06d", n+1)
}

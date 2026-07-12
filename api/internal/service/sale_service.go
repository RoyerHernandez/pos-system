package service

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type saleRepository interface {
	FindAll(startDate, endDate, status string, userID, page, perPage int) ([]domain.Sale, int, error)
	FindByID(id int) (*domain.Sale, error)
	FindDetailsBySaleID(saleID int) ([]domain.SaleDetail, error)
	GetLastSaleCode(tx *sqlx.Tx) (string, error)
	InsertSale(tx *sqlx.Tx, sale *domain.Sale) (int64, error)
	InsertSaleDetail(tx *sqlx.Tx, detail *domain.SaleDetail) error
	UpdateSaleStatus(tx *sqlx.Tx, id int, status string) error
	FindByIDForUpdate(tx *sqlx.Tx, id int) (*domain.Sale, error)
	FindDetailsBySaleIDTx(tx *sqlx.Tx, saleID int) ([]domain.SaleDetail, error)
}

type saleInventoryRepository interface {
	InsertMovement(tx *sqlx.Tx, movement *domain.InventoryMovement) error
}

type SaleRequest struct {
	IDCliente  int
	MetodoPago string
	Productos  []SaleItemRequest
}

type SaleItemRequest struct {
	IDProducto int
	Cantidad   int
}

type SaleResponse struct {
	Sale    domain.Sale
	Details []domain.SaleDetail
}

type SaleService struct {
	db        *sqlx.DB
	sales     saleRepository
	inventory saleInventoryRepository
}

func NewSaleService(db *sqlx.DB, sales saleRepository, inventory saleInventoryRepository) (*SaleService, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	if sales == nil {
		return nil, errors.New("saleRepository must not be nil")
	}
	if inventory == nil {
		return nil, errors.New("inventoryRepository must not be nil")
	}
	return &SaleService{db: db, sales: sales, inventory: inventory}, nil
}

func (s *SaleService) GetAll(startDate, endDate, status string, userID, page, perPage int) ([]domain.Sale, int, error) {
	return s.sales.FindAll(startDate, endDate, status, userID, page, perPage)
}

func (s *SaleService) GetByID(id int) (*SaleResponse, error) {
	sale, err := s.sales.FindByID(id)
	if err != nil {
		return nil, err
	}
	if sale == nil {
		return nil, nil
	}

	details, err := s.sales.FindDetailsBySaleID(id)
	if err != nil {
		return nil, err
	}

	return &SaleResponse{
		Sale:    *sale,
		Details: details,
	}, nil
}

func (s *SaleService) CreateSale(req SaleRequest, userID int) (*SaleResponse, error) {
	if len(req.Productos) == 0 {
		return nil, errors.New("productos is required")
	}
	if req.MetodoPago == "" {
		return nil, errors.New("metodo_pago is required")
	}
	if req.MetodoPago != "Efectivo" && req.MetodoPago != "Tarjeta" && req.MetodoPago != "Transferencia" {
		return nil, errors.New("invalid metodo_pago: must be Efectivo, Tarjeta, or Transferencia")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errors.New("failed to begin transaction")
	}
	defer tx.Rollback()

	// 1. Get open cash register for user
	var caja struct {
		ID int `db:"id"`
	}
	err = tx.Get(&caja, "SELECT id FROM caja WHERE id_usuario = ? AND estado = 'abierta' FOR UPDATE", userID)
	if err != nil {
		return nil, errors.New("no open cash register found for user")
	}

	// 2. Lock and validate stock for each product
	type productInfo struct {
		ID          int     `db:"id"`
		PrecioVenta float64 `db:"precio_venta"`
		Stock       int     `db:"stock"`
	}
	products := make(map[int]*productInfo)
	for _, item := range req.Productos {
		if item.Cantidad <= 0 {
			return nil, fmt.Errorf("cantidad must be greater than 0 for product %d", item.IDProducto)
		}
		var p productInfo
		err = tx.Get(&p, "SELECT id, precio_venta, stock FROM productos WHERE id = ? FOR UPDATE", item.IDProducto)
		if err != nil {
			return nil, fmt.Errorf("product %d not found", item.IDProducto)
		}
		if p.Stock < item.Cantidad {
			return nil, fmt.Errorf("insufficient stock for product %d: available %d, requested %d", item.IDProducto, p.Stock, item.Cantidad)
		}
		products[item.IDProducto] = &p
	}

	// 3. Generate sale code
	lastCode, err := s.sales.GetLastSaleCode(tx)
	if err != nil {
		return nil, errors.New("failed to get last sale code")
	}
	nextNum := 1
	if lastCode != "" {
		parts := strings.Split(lastCode, "-")
		if len(parts) == 2 {
			if n, parseErr := strconv.Atoi(parts[1]); parseErr == nil {
				nextNum = n + 1
			}
		}
	}
	codigoVenta := fmt.Sprintf("VTA-%06d", nextNum)

	// 4. Calculate totals
	var subtotal float64
	for _, item := range req.Productos {
		p := products[item.IDProducto]
		subtotal += float64(item.Cantidad) * p.PrecioVenta
	}
	subtotal = math.Round(subtotal*100) / 100
	impuesto := 0.0
	descuento := 0.0
	total := subtotal + impuesto - descuento

	// 5. Insert sale
	clienteID := &req.IDCliente
	if req.IDCliente == 0 {
		clienteID = nil
	}
	cajaID := caja.ID
	sale := &domain.Sale{
		IDUsuario:   userID,
		IDCliente:   clienteID,
		IDCaja:      &cajaID,
		CodigoVenta: codigoVenta,
		Subtotal:    subtotal,
		Impuesto:    impuesto,
		Descuento:   descuento,
		Total:       total,
		MetodoPago:  req.MetodoPago,
		Estado:      "completada",
	}

	saleID, err := s.sales.InsertSale(tx, sale)
	if err != nil {
		return nil, errors.New("failed to insert sale")
	}
	sale.ID = int(saleID)

	// 6. For each product: insert detail, update stock, insert inventory movement
	var details []domain.SaleDetail
	for _, item := range req.Productos {
		p := products[item.IDProducto]
		lineSubtotal := math.Round(float64(item.Cantidad)*p.PrecioVenta*100) / 100

		detail := &domain.SaleDetail{
			IDVenta:        sale.ID,
			IDProducto:     item.IDProducto,
			Cantidad:       item.Cantidad,
			PrecioUnitario: p.PrecioVenta,
			Descuento:      0,
			Subtotal:       lineSubtotal,
		}
		if err := s.sales.InsertSaleDetail(tx, detail); err != nil {
			return nil, errors.New("failed to insert sale detail")
		}
		details = append(details, *detail)

		// Update product stock and sales count
		_, err = tx.Exec("UPDATE productos SET stock = stock - ?, ventas = ventas + ? WHERE id = ?",
			item.Cantidad, item.Cantidad, item.IDProducto)
		if err != nil {
			return nil, errors.New("failed to update product stock")
		}

		// Insert inventory movement
		obs := fmt.Sprintf("Venta %s", codigoVenta)
		refID := sale.ID
		movement := &domain.InventoryMovement{
			IDProducto:    item.IDProducto,
			IDUsuario:     userID,
			Tipo:          "salida",
			Motivo:        "Venta",
			Cantidad:      item.Cantidad,
			Observaciones: &obs,
			IDReferencia:  &refID,
		}
		if err := s.inventory.InsertMovement(tx, movement); err != nil {
			return nil, errors.New("failed to insert inventory movement")
		}
	}

	// 7. Update client purchases
	if req.IDCliente > 0 {
		_, err = tx.Exec("UPDATE clientes SET total_compras = total_compras + ? WHERE id = ?", total, req.IDCliente)
		if err != nil {
			return nil, errors.New("failed to update client purchases")
		}
	}

	// 8. Update cash register
	updateCaja := "UPDATE caja SET total_ventas = total_ventas + ?"
	switch req.MetodoPago {
	case "Efectivo":
		updateCaja += ", total_efectivo = total_efectivo + ?"
	case "Tarjeta":
		updateCaja += ", total_tarjeta = total_tarjeta + ?"
	case "Transferencia":
		updateCaja += ", total_transferencia = total_transferencia + ?"
	}
	updateCaja += " WHERE id = ?"
	_, err = tx.Exec(updateCaja, total, total, caja.ID)
	if err != nil {
		return nil, errors.New("failed to update cash register")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	return &SaleResponse{
		Sale:    *sale,
		Details: details,
	}, nil
}

func (s *SaleService) CancelSale(saleID int, userID int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return errors.New("failed to begin transaction")
	}
	defer tx.Rollback()

	// 1. Get sale and validate
	sale, err := s.sales.FindByIDForUpdate(tx, saleID)
	if err != nil {
		return errors.New("failed to find sale")
	}
	if sale == nil {
		return errors.New("sale not found")
	}
	if sale.Estado == "cancelada" {
		return errors.New("sale is already cancelled")
	}

	// 2. Get sale details
	details, err := s.sales.FindDetailsBySaleIDTx(tx, saleID)
	if err != nil {
		return errors.New("failed to find sale details")
	}

	// 3. Lock and reverse each product
	for _, detail := range details {
		var lockedID int
		if err := tx.Get(&lockedID, "SELECT id FROM productos WHERE id = ? FOR UPDATE", detail.IDProducto); err != nil {
			return fmt.Errorf("failed to lock product %d", detail.IDProducto)
		}
		_, err = tx.Exec("UPDATE productos SET stock = stock + ?, ventas = ventas - ? WHERE id = ?",
			detail.Cantidad, detail.Cantidad, detail.IDProducto)
		if err != nil {
			return errors.New("failed to restore product stock")
		}

		obs := fmt.Sprintf("Cancelacion venta %s", sale.CodigoVenta)
		refID := saleID
		movement := &domain.InventoryMovement{
			IDProducto:    detail.IDProducto,
			IDUsuario:     userID,
			Tipo:          "entrada",
			Motivo:        "Cancelacion venta",
			Cantidad:      detail.Cantidad,
			Observaciones: &obs,
			IDReferencia:  &refID,
		}
		if err := s.inventory.InsertMovement(tx, movement); err != nil {
			return errors.New("failed to insert reversal movement")
		}
	}

	// 4. Reverse client purchases
	if sale.IDCliente != nil && *sale.IDCliente > 0 {
		_, err = tx.Exec("UPDATE clientes SET total_compras = total_compras - ? WHERE id = ?", sale.Total, *sale.IDCliente)
		if err != nil {
			return errors.New("failed to reverse client purchases")
		}
	}

	// 5. Update sale status
	if err := s.sales.UpdateSaleStatus(tx, saleID, "cancelada"); err != nil {
		return errors.New("failed to update sale status")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("failed to commit transaction")
	}

	return nil
}

package services

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
	"github.com/jmoiron/sqlx"
)

type SaleService struct {
	db            *sqlx.DB
	saleRepo      *repositories.SaleRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewSaleService(db *sqlx.DB, saleRepo *repositories.SaleRepository, inventoryRepo *repositories.InventoryRepository) *SaleService {
	return &SaleService{
		db:            db,
		saleRepo:      saleRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *SaleService) GetAll(startDate, endDate, status string, userID, page, perPage int) ([]models.Sale, int, error) {
	return s.saleRepo.FindAll(startDate, endDate, status, userID, page, perPage)
}

func (s *SaleService) GetByID(id int) (*models.SaleResponse, error) {
	sale, err := s.saleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if sale == nil {
		return nil, nil
	}

	details, err := s.saleRepo.FindDetailsBySaleID(id)
	if err != nil {
		return nil, err
	}

	return &models.SaleResponse{
		Sale:    *sale,
		Details: details,
	}, nil
}

func (s *SaleService) CreateSale(req models.CreateSaleRequest, userID int) (*models.SaleResponse, error) {
	if len(req.Productos) == 0 {
		return nil, fmt.Errorf("productos is required")
	}
	if req.MetodoPago == "" {
		return nil, fmt.Errorf("metodo_pago is required")
	}
	if req.MetodoPago != "Efectivo" && req.MetodoPago != "Tarjeta" && req.MetodoPago != "Transferencia" {
		return nil, fmt.Errorf("invalid metodo_pago: must be Efectivo, Tarjeta, or Transferencia")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Get open cash register for user
	var caja struct {
		ID int `db:"id"`
	}
	err = tx.Get(&caja, "SELECT id FROM caja WHERE id_usuario = ? AND estado = 'abierta' FOR UPDATE", userID)
	if err != nil {
		return nil, fmt.Errorf("no open cash register found for user")
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
	lastCode, err := s.saleRepo.GetLastSaleCode(tx)
	if err != nil {
		return nil, fmt.Errorf("get last sale code: %w", err)
	}
	nextNum := 1
	if lastCode != "" {
		parts := strings.Split(lastCode, "-")
		if len(parts) == 2 {
			if n, err := strconv.Atoi(parts[1]); err == nil {
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
	sale := &models.Sale{
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

	saleID, err := s.saleRepo.InsertSale(tx, sale)
	if err != nil {
		return nil, fmt.Errorf("insert sale: %w", err)
	}
	sale.ID = int(saleID)

	// 6. For each product: insert detail, update stock, insert inventory movement
	var details []models.SaleDetail
	for _, item := range req.Productos {
		p := products[item.IDProducto]
		lineSubtotal := math.Round(float64(item.Cantidad)*p.PrecioVenta*100) / 100

		detail := &models.SaleDetail{
			IDVenta:        sale.ID,
			IDProducto:     item.IDProducto,
			Cantidad:       item.Cantidad,
			PrecioUnitario: p.PrecioVenta,
			Descuento:      0,
			Subtotal:       lineSubtotal,
		}
		if err := s.saleRepo.InsertSaleDetail(tx, detail); err != nil {
			return nil, fmt.Errorf("insert detail: %w", err)
		}
		details = append(details, *detail)

		// Update product stock and sales count
		_, err = tx.Exec("UPDATE productos SET stock = stock - ?, ventas = ventas + ? WHERE id = ?",
			item.Cantidad, item.Cantidad, item.IDProducto)
		if err != nil {
			return nil, fmt.Errorf("update product stock: %w", err)
		}

		// Insert inventory movement
		obs := fmt.Sprintf("Venta %s", codigoVenta)
		refID := sale.ID
		movement := &models.InventoryMovement{
			IDProducto:    item.IDProducto,
			IDUsuario:     userID,
			Tipo:          "salida",
			Motivo:        "Venta",
			Cantidad:      item.Cantidad,
			Observaciones: &obs,
			IDReferencia:  &refID,
		}
		if err := s.inventoryRepo.InsertMovement(tx, movement); err != nil {
			return nil, fmt.Errorf("insert inventory movement: %w", err)
		}
	}

	// 7. Update client purchases
	if req.IDCliente > 0 {
		_, err = tx.Exec("UPDATE clientes SET compras = compras + ? WHERE id = ?", total, req.IDCliente)
		if err != nil {
			return nil, fmt.Errorf("update client purchases: %w", err)
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
		return nil, fmt.Errorf("update cash register: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &models.SaleResponse{
		Sale:    *sale,
		Details: details,
	}, nil
}

func (s *SaleService) CancelSale(saleID int, userID int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Get sale and validate
	sale, err := s.saleRepo.FindByIDForUpdate(tx, saleID)
	if err != nil {
		return fmt.Errorf("find sale: %w", err)
	}
	if sale == nil {
		return fmt.Errorf("sale not found")
	}
	if sale.Estado == "cancelada" {
		return fmt.Errorf("sale is already cancelled")
	}

	// 2. Get sale details
	details, err := s.saleRepo.FindDetailsBySaleIDTx(tx, saleID)
	if err != nil {
		return fmt.Errorf("find sale details: %w", err)
	}

	// 3. Lock and reverse each product
	for _, detail := range details {
		var lockedID int
		if err := tx.Get(&lockedID, "SELECT id FROM productos WHERE id = ? FOR UPDATE", detail.IDProducto); err != nil {
			return fmt.Errorf("lock product %d: %w", detail.IDProducto, err)
		}
		_, err = tx.Exec("UPDATE productos SET stock = stock + ?, ventas = ventas - ? WHERE id = ?",
			detail.Cantidad, detail.Cantidad, detail.IDProducto)
		if err != nil {
			return fmt.Errorf("restore product stock: %w", err)
		}

		obs := fmt.Sprintf("Cancelacion venta %s", sale.CodigoVenta)
		refID := saleID
		movement := &models.InventoryMovement{
			IDProducto:    detail.IDProducto,
			IDUsuario:     userID,
			Tipo:          "entrada",
			Motivo:        "Cancelacion venta",
			Cantidad:      detail.Cantidad,
			Observaciones: &obs,
			IDReferencia:  &refID,
		}
		if err := s.inventoryRepo.InsertMovement(tx, movement); err != nil {
			return fmt.Errorf("insert reversal movement: %w", err)
		}
	}

	// 4. Reverse client purchases
	if sale.IDCliente != nil && *sale.IDCliente > 0 {
		_, err = tx.Exec("UPDATE clientes SET compras = compras - ? WHERE id = ?", sale.Total, *sale.IDCliente)
		if err != nil {
			return fmt.Errorf("reverse client purchases: %w", err)
		}
	}

	// 5. Update sale status
	if err := s.saleRepo.UpdateSaleStatus(tx, saleID, "cancelada"); err != nil {
		return fmt.Errorf("update sale status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

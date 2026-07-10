package services

import (
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
	"github.com/jmoiron/sqlx"
)

type InventoryService struct {
	db            *sqlx.DB
	inventoryRepo *repositories.InventoryRepository
}

func NewInventoryService(db *sqlx.DB, inventoryRepo *repositories.InventoryRepository) *InventoryService {
	return &InventoryService{
		db:            db,
		inventoryRepo: inventoryRepo,
	}
}

func (s *InventoryService) GetMovements(productID int, tipo, startDate, endDate string, page, perPage int) ([]models.InventoryMovement, int, error) {
	return s.inventoryRepo.FindAll(productID, tipo, startDate, endDate, page, perPage)
}

func (s *InventoryService) CreateManualMovement(req models.CreateMovementRequest, userID int) (*models.InventoryMovement, error) {
	if req.IDProducto <= 0 {
		return nil, fmt.Errorf("id_producto is required")
	}
	if req.Cantidad <= 0 {
		return nil, fmt.Errorf("cantidad must be greater than 0")
	}
	if req.Tipo != "entrada" && req.Tipo != "salida" {
		return nil, fmt.Errorf("tipo must be 'entrada' or 'salida'")
	}
	if req.Motivo == "" {
		return nil, fmt.Errorf("motivo is required")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Validate product exists and lock row
	var product struct {
		ID    int `db:"id"`
		Stock int `db:"stock"`
	}
	err = tx.Get(&product, "SELECT id, stock FROM productos WHERE id = ? FOR UPDATE", req.IDProducto)
	if err != nil {
		return nil, fmt.Errorf("product %d not found", req.IDProducto)
	}

	// For salida, validate sufficient stock
	if req.Tipo == "salida" && product.Stock < req.Cantidad {
		return nil, fmt.Errorf("insufficient stock: available %d, requested %d", product.Stock, req.Cantidad)
	}

	// Insert movement
	movement := &models.InventoryMovement{
		IDProducto:    req.IDProducto,
		IDUsuario:     userID,
		Tipo:          req.Tipo,
		Motivo:        req.Motivo,
		Cantidad:      req.Cantidad,
		Observaciones: req.Observaciones,
	}

	if err := s.inventoryRepo.InsertMovement(tx, movement); err != nil {
		return nil, fmt.Errorf("insert movement: %w", err)
	}

	// Update product stock
	var stockUpdate string
	if req.Tipo == "entrada" {
		stockUpdate = "UPDATE productos SET stock = stock + ? WHERE id = ?"
	} else {
		stockUpdate = "UPDATE productos SET stock = stock - ? WHERE id = ?"
	}
	if _, err := tx.Exec(stockUpdate, req.Cantidad, req.IDProducto); err != nil {
		return nil, fmt.Errorf("update product stock: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return movement, nil
}

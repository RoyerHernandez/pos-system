package service

import (
	"errors"
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type inventoryRepository interface {
	FindAll(productID int, tipo, startDate, endDate string, page, perPage int) ([]domain.InventoryMovement, int, error)
	InsertMovement(tx *sqlx.Tx, movement *domain.InventoryMovement) error
}

type CreateMovementRequest struct {
	IDProducto    int
	Tipo          string
	Motivo        string
	Cantidad      int
	Observaciones *string
}

type InventoryService struct {
	db   *sqlx.DB
	repo inventoryRepository
}

func NewInventoryService(repo inventoryRepository, db *sqlx.DB) (*InventoryService, error) {
	if repo == nil {
		return nil, errors.New("inventoryRepository must not be nil")
	}
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &InventoryService{db: db, repo: repo}, nil
}

func (s *InventoryService) GetMovements(productID int, tipo, startDate, endDate string, page, perPage int) ([]domain.InventoryMovement, int, error) {
	return s.repo.FindAll(productID, tipo, startDate, endDate, page, perPage)
}

func (s *InventoryService) CreateManualMovement(req CreateMovementRequest, userID int) (*domain.InventoryMovement, error) {
	if req.IDProducto <= 0 {
		return nil, errors.New("id_producto is required")
	}
	if req.Cantidad <= 0 {
		return nil, errors.New("cantidad must be greater than 0")
	}
	if req.Tipo != "entrada" && req.Tipo != "salida" {
		return nil, errors.New("tipo must be 'entrada' or 'salida'")
	}
	if req.Motivo == "" {
		return nil, errors.New("motivo is required")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errors.New("failed to begin transaction")
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
	movement := &domain.InventoryMovement{
		IDProducto:    req.IDProducto,
		IDUsuario:     userID,
		Tipo:          req.Tipo,
		Motivo:        req.Motivo,
		Cantidad:      req.Cantidad,
		Observaciones: req.Observaciones,
	}

	if err := s.repo.InsertMovement(tx, movement); err != nil {
		return nil, errors.New("failed to insert movement")
	}

	// Update product stock
	var stockUpdate string
	if req.Tipo == "entrada" {
		stockUpdate = "UPDATE productos SET stock = stock + ? WHERE id = ?"
	} else {
		stockUpdate = "UPDATE productos SET stock = stock - ? WHERE id = ?"
	}
	if _, err := tx.Exec(stockUpdate, req.Cantidad, req.IDProducto); err != nil {
		return nil, errors.New("failed to update product stock")
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	return movement, nil
}

package repositories

import (
	"fmt"
	"time"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/jmoiron/sqlx"
)

type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) InsertMovement(tx *sqlx.Tx, movement *models.InventoryMovement) error {
	result, err := tx.Exec(
		`INSERT INTO movimientos_inventario (id_producto, id_usuario, tipo, motivo, cantidad, observaciones, id_referencia)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		movement.IDProducto, movement.IDUsuario, movement.Tipo,
		movement.Motivo, movement.Cantidad, movement.Observaciones,
		movement.IDReferencia,
	)
	if err != nil {
		return fmt.Errorf("insert inventory movement: %w", err)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	movement.ID = int(lastID)
	now := time.Now()
	movement.Fecha = now
	return nil
}

func (r *InventoryRepository) FindAll(productID int, tipo, startDate, endDate string, page, perPage int) ([]models.InventoryMovement, int, error) {
	query := "SELECT * FROM movimientos_inventario WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM movimientos_inventario WHERE 1=1"
	var args []interface{}

	if productID > 0 {
		query += " AND id_producto = ?"
		countQuery += " AND id_producto = ?"
		args = append(args, productID)
	}
	if tipo != "" {
		query += " AND tipo = ?"
		countQuery += " AND tipo = ?"
		args = append(args, tipo)
	}
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

	var total int
	if err := r.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("count inventory movements: %w", err)
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

	var movements []models.InventoryMovement
	if err := r.db.Select(&movements, query, args...); err != nil {
		return nil, 0, fmt.Errorf("find inventory movements: %w", err)
	}
	return movements, total, nil
}

package repository

import (
	"errors"
	"log"
	"time"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) (*InventoryRepository, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &InventoryRepository{db: db}, nil
}

func (r *InventoryRepository) InsertMovement(tx *sqlx.Tx, movement *domain.InventoryMovement) error {
	dao := toInventoryMovementDAO(*movement)
	result, err := tx.Exec(
		`INSERT INTO movimientos_inventario (id_producto, id_usuario, tipo, motivo, cantidad, observaciones, id_referencia)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		dao.IDProducto, dao.IDUsuario, dao.Tipo,
		dao.Motivo, dao.Cantidad, dao.Observaciones,
		dao.IDReferencia,
	)
	if err != nil {
		log.Printf("repository - InsertMovement: %v", err)
		return errors.New("failed to insert inventory movement")
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Printf("repository - InsertMovement last insert id: %v", err)
		return errors.New("failed to get last insert id for inventory movement")
	}
	movement.ID = int(lastID)
	movement.Fecha = time.Now()
	return nil
}

func (r *InventoryRepository) FindAll(productID int, tipo, startDate, endDate string, page, perPage int) ([]domain.InventoryMovement, int, error) {
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
		log.Printf("repository - FindAll(inventory) count: %v", err)
		return nil, 0, errors.New("failed to count inventory movements")
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

	var daos []inventoryMovementDAO
	if err := r.db.Select(&daos, query, args...); err != nil {
		log.Printf("repository - FindAll(inventory): %v", err)
		return nil, 0, errors.New("failed to find inventory movements")
	}
	movements := make([]domain.InventoryMovement, len(daos))
	for i := range daos {
		movements[i] = daos[i].toDomain()
	}
	return movements, total, nil
}

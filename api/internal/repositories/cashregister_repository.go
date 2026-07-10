package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/jmoiron/sqlx"
)

type CashRegisterRepository struct {
	db *sqlx.DB
}

func NewCashRegisterRepository(db *sqlx.DB) *CashRegisterRepository {
	return &CashRegisterRepository{db: db}
}

func (r *CashRegisterRepository) FindAll(page, perPage int) ([]models.CashRegister, int, error) {
	var total int
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM caja"); err != nil {
		return nil, 0, fmt.Errorf("count cash registers: %w", err)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 25
	}
	offset := (page - 1) * perPage

	var registers []models.CashRegister
	err := r.db.Select(&registers, "SELECT * FROM caja ORDER BY id DESC LIMIT ? OFFSET ?", perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("find cash registers: %w", err)
	}
	return registers, total, nil
}

func (r *CashRegisterRepository) FindByID(id int) (*models.CashRegister, error) {
	var reg models.CashRegister
	err := r.db.Get(&reg, "SELECT * FROM caja WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find cash register: %w", err)
	}
	return &reg, nil
}

func (r *CashRegisterRepository) FindOpenByUserID(userID int) (*models.CashRegister, error) {
	var reg models.CashRegister
	err := r.db.Get(&reg, "SELECT * FROM caja WHERE id_usuario = ? AND estado = 'abierta'", userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find open cash register: %w", err)
	}
	return &reg, nil
}

func (r *CashRegisterRepository) Insert(reg *models.CashRegister) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO caja (id_usuario, monto_apertura) VALUES (?, ?)",
		reg.IDUsuario, reg.MontoApertura,
	)
	if err != nil {
		return 0, fmt.Errorf("insert cash register: %w", err)
	}
	return result.LastInsertId()
}

func (r *CashRegisterRepository) FindOpenByUserIDTx(tx *sqlx.Tx, userID int) (*models.CashRegister, error) {
	var reg models.CashRegister
	err := tx.Get(&reg, "SELECT * FROM caja WHERE id_usuario = ? AND estado = 'abierta' FOR UPDATE", userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find open cash register for update: %w", err)
	}
	return &reg, nil
}

func (r *CashRegisterRepository) InsertTx(tx *sqlx.Tx, reg *models.CashRegister) (int64, error) {
	result, err := tx.Exec(
		"INSERT INTO caja (id_usuario, monto_apertura) VALUES (?, ?)",
		reg.IDUsuario, reg.MontoApertura,
	)
	if err != nil {
		return 0, fmt.Errorf("insert cash register: %w", err)
	}
	return result.LastInsertId()
}

func (r *CashRegisterRepository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *CashRegisterRepository) Close(id int, montoCierre float64) error {
	_, err := r.db.Exec(
		"UPDATE caja SET estado = 'cerrada', monto_cierre = ?, fecha_cierre = NOW() WHERE id = ?",
		montoCierre, id,
	)
	if err != nil {
		return fmt.Errorf("close cash register: %w", err)
	}
	return nil
}

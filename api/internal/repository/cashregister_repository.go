package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

type CashRegisterRepository struct {
	db *sqlx.DB
}

func NewCashRegisterRepository(db *sqlx.DB) (*CashRegisterRepository, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &CashRegisterRepository{db: db}, nil
}

func (r *CashRegisterRepository) FindAll(page, perPage int) ([]domain.CashRegister, int, error) {
	var total int
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM caja"); err != nil {
		log.Printf("repository - FindAll(cashregisters) count: %v", err)
		return nil, 0, errors.New("failed to count cash registers")
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 25
	}
	offset := (page - 1) * perPage

	var daos []cashRegisterDAO
	err := r.db.Select(&daos, "SELECT * FROM caja ORDER BY id DESC LIMIT ? OFFSET ?", perPage, offset)
	if err != nil {
		log.Printf("repository - FindAll(cashregisters): %v", err)
		return nil, 0, errors.New("failed to find cash registers")
	}
	registers := make([]domain.CashRegister, len(daos))
	for i := range daos {
		registers[i] = daos[i].toDomain()
	}
	return registers, total, nil
}

func (r *CashRegisterRepository) FindByID(id int) (*domain.CashRegister, error) {
	var dao cashRegisterDAO
	err := r.db.Get(&dao, "SELECT * FROM caja WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFoundError{Message: "cash register not found"}
	}
	if err != nil {
		log.Printf("repository - FindByID(cashregister): %v", err)
		return nil, errors.New("failed to find cash register")
	}
	cr := dao.toDomain()
	return &cr, nil
}

func (r *CashRegisterRepository) FindOpenByUserID(userID int) (*domain.CashRegister, error) {
	var dao cashRegisterDAO
	err := r.db.Get(&dao, "SELECT * FROM caja WHERE id_usuario = ? AND estado = 'abierta'", userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		log.Printf("repository - FindOpenByUserID: %v", err)
		return nil, errors.New("failed to find open cash register")
	}
	cr := dao.toDomain()
	return &cr, nil
}

func (r *CashRegisterRepository) Insert(reg *domain.CashRegister) (int64, error) {
	dao := toCashRegisterDAO(*reg)
	result, err := r.db.Exec(
		"INSERT INTO caja (id_usuario, monto_apertura) VALUES (?, ?)",
		dao.IDUsuario, dao.MontoApertura,
	)
	if err != nil {
		log.Printf("repository - Insert(cashregister): %v", err)
		return 0, errors.New("failed to insert cash register")
	}
	return result.LastInsertId()
}

func (r *CashRegisterRepository) FindOpenByUserIDTx(tx *sqlx.Tx, userID int) (*domain.CashRegister, error) {
	var dao cashRegisterDAO
	err := tx.Get(&dao, "SELECT * FROM caja WHERE id_usuario = ? AND estado = 'abierta' FOR UPDATE", userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		log.Printf("repository - FindOpenByUserIDTx: %v", err)
		return nil, errors.New("failed to find open cash register for update")
	}
	cr := dao.toDomain()
	return &cr, nil
}

func (r *CashRegisterRepository) InsertTx(tx *sqlx.Tx, reg *domain.CashRegister) (int64, error) {
	dao := toCashRegisterDAO(*reg)
	result, err := tx.Exec(
		"INSERT INTO caja (id_usuario, monto_apertura) VALUES (?, ?)",
		dao.IDUsuario, dao.MontoApertura,
	)
	if err != nil {
		log.Printf("repository - InsertTx(cashregister): %v", err)
		return 0, errors.New("failed to insert cash register in transaction")
	}
	return result.LastInsertId()
}

func (r *CashRegisterRepository) BeginTx() (*sqlx.Tx, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("repository - BeginTx: %v", err)
		return nil, errors.New("failed to begin transaction")
	}
	return tx, nil
}

func (r *CashRegisterRepository) Close(id int, montoCierre float64) error {
	_, err := r.db.Exec(
		"UPDATE caja SET estado = 'cerrada', monto_cierre = ?, fecha_cierre = NOW() WHERE id = ?",
		montoCierre, id,
	)
	if err != nil {
		log.Printf("repository - Close(cashregister): %v", err)
		return errors.New("failed to close cash register")
	}
	return nil
}

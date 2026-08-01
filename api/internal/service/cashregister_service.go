package service

import (
	"errors"
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type cashRegisterRepository interface {
	FindAll(page, perPage int) ([]domain.CashRegister, int, error)
	FindByID(id int) (*domain.CashRegister, error)
	FindOpenByUserID(userID int) (*domain.CashRegister, error)
	FindOpenByUserIDTx(tx *sqlx.Tx, userID int) (*domain.CashRegister, error)
	Insert(reg *domain.CashRegister) (int64, error)
	InsertTx(tx *sqlx.Tx, reg *domain.CashRegister) (int64, error)
	BeginTx() (*sqlx.Tx, error)
	Close(id int, montoCierre float64) error
}

type CashRegisterService struct {
	repo cashRegisterRepository
	db   *sqlx.DB
}

func NewCashRegisterService(repo cashRegisterRepository, db *sqlx.DB) (*CashRegisterService, error) {
	if repo == nil {
		return nil, errors.New("cashRegisterRepository must not be nil")
	}
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &CashRegisterService{repo: repo, db: db}, nil
}

func (s *CashRegisterService) OpenRegister(userID int, montoApertura float64) (*domain.CashRegister, error) {
	if montoApertura < 0 {
		return nil, errors.New("monto_apertura must be >= 0")
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return nil, errors.New("failed to begin transaction")
	}
	defer tx.Rollback()

	// Check-and-insert within a transaction using SELECT ... FOR UPDATE
	existing, err := s.repo.FindOpenByUserIDTx(tx, userID)
	if err != nil {
		return nil, errors.New("failed to check existing register")
	}
	if existing != nil {
		return nil, fmt.Errorf("user already has an open cash register (id: %d)", existing.ID)
	}

	reg := &domain.CashRegister{
		IDUsuario:     userID,
		MontoApertura: montoApertura,
	}

	id, err := s.repo.InsertTx(tx, reg)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.New("failed to commit transaction")
	}

	reg.ID = int(id)
	reg.Estado = "abierta"

	return reg, nil
}

func (s *CashRegisterService) CloseRegister(registerID, userID int) (*domain.CashRegister, error) {
	reg, err := s.repo.FindByID(registerID)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		return nil, errors.New("cash register not found")
	}
	if reg.Estado != "abierta" {
		return nil, errors.New("cash register is already closed")
	}
	if reg.IDUsuario != userID {
		return nil, errors.New("cash register does not belong to this user")
	}

	// monto_cierre represents the physical cash expected in the drawer:
	// opening amount plus cash sales only. Card and transfer payments are
	// settled externally and do not affect the physical cash count.
	montoCierre := reg.MontoApertura + reg.TotalEfectivo
	if err := s.repo.Close(registerID, montoCierre); err != nil {
		return nil, err
	}

	// Return updated register
	return s.repo.FindByID(registerID)
}

func (s *CashRegisterService) GetCurrentOpen(userID int) (*domain.CashRegister, error) {
	return s.repo.FindOpenByUserID(userID)
}

func (s *CashRegisterService) GetAll(page, perPage int) ([]domain.CashRegister, int, error) {
	return s.repo.FindAll(page, perPage)
}

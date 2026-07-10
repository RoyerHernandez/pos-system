package services

import (
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
)

type CashRegisterService struct {
	cashRepo *repositories.CashRegisterRepository
}

func NewCashRegisterService(cashRepo *repositories.CashRegisterRepository) *CashRegisterService {
	return &CashRegisterService{cashRepo: cashRepo}
}

func (s *CashRegisterService) OpenRegister(userID int, montoApertura float64) (*models.CashRegister, error) {
	if montoApertura < 0 {
		return nil, fmt.Errorf("monto_apertura must be >= 0")
	}

	tx, err := s.cashRepo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check-and-insert within a transaction using SELECT ... FOR UPDATE
	existing, err := s.cashRepo.FindOpenByUserIDTx(tx, userID)
	if err != nil {
		return nil, fmt.Errorf("check existing register: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user already has an open cash register (id: %d)", existing.ID)
	}

	reg := &models.CashRegister{
		IDUsuario:     userID,
		MontoApertura: montoApertura,
	}

	id, err := s.cashRepo.InsertTx(tx, reg)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	reg.ID = int(id)
	reg.Estado = "abierta"

	return reg, nil
}

func (s *CashRegisterService) CloseRegister(registerID, userID int) (*models.CashRegister, error) {
	reg, err := s.cashRepo.FindByID(registerID)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		return nil, fmt.Errorf("cash register not found")
	}
	if reg.Estado != "abierta" {
		return nil, fmt.Errorf("cash register is already closed")
	}
	if reg.IDUsuario != userID {
		return nil, fmt.Errorf("cash register does not belong to this user")
	}

	// monto_cierre represents the physical cash expected in the drawer:
	// opening amount plus cash sales only. Card and transfer payments are
	// settled externally and do not affect the physical cash count.
	montoCierre := reg.MontoApertura + reg.TotalEfectivo
	if err := s.cashRepo.Close(registerID, montoCierre); err != nil {
		return nil, err
	}

	// Return updated register
	return s.cashRepo.FindByID(registerID)
}

func (s *CashRegisterService) GetCurrentOpen(userID int) (*models.CashRegister, error) {
	return s.cashRepo.FindOpenByUserID(userID)
}

func (s *CashRegisterService) GetAll(page, perPage int) ([]models.CashRegister, int, error) {
	return s.cashRepo.FindAll(page, perPage)
}

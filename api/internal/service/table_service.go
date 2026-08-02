package service

import (
	"errors"
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
)

type tableRepository interface {
	FindAll() ([]domain.Table, error)
	FindByID(id int) (*domain.Table, error)
	Create(t *domain.Table) (int64, error)
	Update(t *domain.Table) error
	Delete(id int) error
	ExistsByNumero(numero int, excludeID int) (bool, error)
}

type TableService struct {
	tables tableRepository
}

func NewTableService(tables tableRepository) (*TableService, error) {
	if tables == nil {
		return nil, errors.New("tableRepository must not be nil")
	}
	return &TableService{tables: tables}, nil
}

func (s *TableService) GetAll() ([]domain.Table, error) {
	return s.tables.FindAll()
}

func (s *TableService) GetByID(id int) (*domain.Table, error) {
	return s.tables.FindByID(id)
}

func (s *TableService) Create(t *domain.Table) (int64, error) {
	if t.Numero <= 0 {
		return 0, apperror.ValidationError{Message: "numero must be greater than 0"}
	}
	if t.Capacidad <= 0 {
		t.Capacidad = 4
	}

	exists, err := s.tables.ExistsByNumero(t.Numero, 0)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, apperror.ConflictError{Message: fmt.Sprintf("table number %d already exists", t.Numero)}
	}

	return s.tables.Create(t)
}

func (s *TableService) Update(t *domain.Table) error {
	if t.Numero <= 0 {
		return apperror.ValidationError{Message: "numero must be greater than 0"}
	}

	exists, err := s.tables.ExistsByNumero(t.Numero, t.ID)
	if err != nil {
		return err
	}
	if exists {
		return apperror.ConflictError{Message: fmt.Sprintf("table number %d already exists", t.Numero)}
	}

	return s.tables.Update(t)
}

func (s *TableService) Delete(id int) error {
	return s.tables.Delete(id)
}

package service

import (
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
)

type clientRepository interface {
	FindAll() ([]domain.Client, error)
	FindByID(id int) (*domain.Client, error)
	Create(nombre string, documento, email, telefono, direccion, fechaNacimiento *string) (int64, error)
	Update(id int, nombre string, documento, email, telefono, direccion, fechaNacimiento *string) error
	Delete(id int) error
}

type ClientService struct {
	clients clientRepository
}

func NewClientService(clients clientRepository) (*ClientService, error) {
	if clients == nil {
		return nil, errors.New("clientRepository must not be nil")
	}
	return &ClientService{clients: clients}, nil
}

func (s *ClientService) GetAll() ([]domain.Client, error) {
	return s.clients.FindAll()
}

func (s *ClientService) GetByID(id int) (*domain.Client, error) {
	return s.clients.FindByID(id)
}

func (s *ClientService) Create(nombre string, documento, email, telefono, direccion, fechaNacimiento *string) (int64, error) {
	if nombre == "" {
		return 0, errors.New("nombre is required")
	}
	return s.clients.Create(nombre, documento, email, telefono, direccion, fechaNacimiento)
}

func (s *ClientService) Update(id int, nombre string, documento, email, telefono, direccion, fechaNacimiento *string) error {
	if nombre == "" {
		return errors.New("nombre is required")
	}
	return s.clients.Update(id, nombre, documento, email, telefono, direccion, fechaNacimiento)
}

func (s *ClientService) Delete(id int) error {
	return s.clients.Delete(id)
}

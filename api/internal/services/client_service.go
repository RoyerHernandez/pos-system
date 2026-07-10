package services

import (
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
)

type ClientService struct {
	clientRepo *repositories.ClientRepository
}

func NewClientService(clientRepo *repositories.ClientRepository) *ClientService {
	return &ClientService{clientRepo: clientRepo}
}

func (s *ClientService) GetAll() ([]models.Client, error) {
	return s.clientRepo.FindAll()
}

func (s *ClientService) GetByID(id int) (*models.Client, error) {
	return s.clientRepo.FindByID(id)
}

func (s *ClientService) Create(nombre string, documento, email, telefono, direccion, fechaNacimiento *string) (int64, error) {
	if nombre == "" {
		return 0, errors.New("nombre is required")
	}
	return s.clientRepo.Create(nombre, documento, email, telefono, direccion, fechaNacimiento)
}

func (s *ClientService) Update(id int, nombre string, documento, email, telefono, direccion, fechaNacimiento *string) error {
	if nombre == "" {
		return errors.New("nombre is required")
	}
	return s.clientRepo.Update(id, nombre, documento, email, telefono, direccion, fechaNacimiento)
}

func (s *ClientService) Delete(id int) error {
	return s.clientRepo.Delete(id)
}

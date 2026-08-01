package service

import (
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
)

type categoryRepository interface {
	FindAll() ([]domain.Category, error)
	FindByID(id int) (*domain.Category, error)
	Create(nombre, descripcion string) (int64, error)
	Update(id int, nombre, descripcion string) error
	Delete(id int) error
}

type CategoryService struct {
	categories categoryRepository
}

func NewCategoryService(categories categoryRepository) (*CategoryService, error) {
	if categories == nil {
		return nil, errors.New("categoryRepository must not be nil")
	}
	return &CategoryService{categories: categories}, nil
}

func (s *CategoryService) GetAll() ([]domain.Category, error) {
	return s.categories.FindAll()
}

func (s *CategoryService) GetByID(id int) (*domain.Category, error) {
	return s.categories.FindByID(id)
}

func (s *CategoryService) Create(nombre, descripcion string) (int64, error) {
	if nombre == "" {
		return 0, errors.New("nombre is required")
	}
	return s.categories.Create(nombre, descripcion)
}

func (s *CategoryService) Update(id int, nombre, descripcion string) error {
	if nombre == "" {
		return errors.New("nombre is required")
	}
	return s.categories.Update(id, nombre, descripcion)
}

func (s *CategoryService) Delete(id int) error {
	return s.categories.Delete(id)
}

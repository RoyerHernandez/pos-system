package services

import (
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *CategoryService) Create(nombre, descripcion string) (int64, error) {
	if nombre == "" {
		return 0, errors.New("nombre is required")
	}
	return s.categoryRepo.Create(nombre, descripcion)
}

func (s *CategoryService) Update(id int, nombre, descripcion string) error {
	if nombre == "" {
		return errors.New("nombre is required")
	}
	return s.categoryRepo.Update(id, nombre, descripcion)
}

func (s *CategoryService) Delete(id int) error {
	return s.categoryRepo.Delete(id)
}

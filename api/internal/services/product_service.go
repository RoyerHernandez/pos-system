package services

import (
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) GetAll(categoryID int) ([]models.Product, error) {
	return s.productRepo.FindAll(categoryID)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) Create(idCategoria int, descripcion string, stock int, precioCompra, precioVenta float64) (int64, error) {
	if descripcion == "" {
		return 0, errors.New("descripcion is required")
	}
	if precioVenta <= 0 {
		return 0, errors.New("precio_venta must be greater than 0")
	}
	return s.productRepo.Create(idCategoria, descripcion, stock, precioCompra, precioVenta)
}

func (s *ProductService) Update(id int, idCategoria int, descripcion string, precioCompra, precioVenta float64) error {
	if descripcion == "" {
		return errors.New("descripcion is required")
	}
	if precioVenta <= 0 {
		return errors.New("precio_venta must be greater than 0")
	}
	return s.productRepo.Update(id, idCategoria, descripcion, precioCompra, precioVenta)
}

func (s *ProductService) UpdateImage(id int, imagen string) error {
	return s.productRepo.UpdateImage(id, imagen)
}

func (s *ProductService) Delete(id int) error {
	return s.productRepo.Delete(id)
}

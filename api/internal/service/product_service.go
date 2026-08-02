package service

import (
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
)

type productRepository interface {
	FindAll(categoryID int) ([]domain.Product, error)
	FindByID(id int) (*domain.Product, error)
	Create(codigo string, idCategoria int, descripcion string, stock int, precioCompra, precioVenta float64) (int64, error)
	Update(id int, idCategoria int, descripcion string, precioCompra, precioVenta float64) error
	UpdateImage(id int, imagen string) error
	Delete(id int) error
}

type ProductService struct {
	products productRepository
}

func NewProductService(products productRepository) (*ProductService, error) {
	if products == nil {
		return nil, errors.New("productRepository must not be nil")
	}
	return &ProductService{products: products}, nil
}

func (s *ProductService) GetAll(categoryID int) ([]domain.Product, error) {
	return s.products.FindAll(categoryID)
}

func (s *ProductService) GetByID(id int) (*domain.Product, error) {
	return s.products.FindByID(id)
}

func (s *ProductService) Create(codigo string, idCategoria int, descripcion string, stock int, precioCompra, precioVenta float64) (int64, error) {
	if codigo == "" {
		return 0, errors.New("codigo is required")
	}
	if descripcion == "" {
		return 0, errors.New("descripcion is required")
	}
	if precioVenta <= 0 {
		return 0, errors.New("precio_venta must be greater than 0")
	}
	return s.products.Create(codigo, idCategoria, descripcion, stock, precioCompra, precioVenta)
}

func (s *ProductService) Update(id int, idCategoria int, descripcion string, precioCompra, precioVenta float64) error {
	if descripcion == "" {
		return errors.New("descripcion is required")
	}
	if precioVenta <= 0 {
		return errors.New("precio_venta must be greater than 0")
	}
	return s.products.Update(id, idCategoria, descripcion, precioCompra, precioVenta)
}

func (s *ProductService) UpdateImage(id int, imagen string) error {
	return s.products.UpdateImage(id, imagen)
}

func (s *ProductService) Delete(id int) error {
	return s.products.Delete(id)
}

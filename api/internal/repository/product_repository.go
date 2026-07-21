package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) (*ProductRepository, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &ProductRepository{db: db}, nil
}

func (r *ProductRepository) FindAll(categoryID int) ([]domain.Product, error) {
	var daos []productDAO
	var err error
	if categoryID > 0 {
		err = r.db.Select(&daos, "SELECT * FROM productos WHERE id_categoria = ?", categoryID)
	} else {
		err = r.db.Select(&daos, "SELECT * FROM productos")
	}
	if err != nil {
		log.Printf("repository - FindAll(products): %v", err)
		return nil, errors.New("failed to find products")
	}
	products := make([]domain.Product, len(daos))
	for i := range daos {
		products[i] = daos[i].toDomain()
	}
	return products, nil
}

func (r *ProductRepository) FindByID(id int) (*domain.Product, error) {
	var dao productDAO
	err := r.db.Get(&dao, "SELECT * FROM productos WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFoundError{Message: "product not found"}
	}
	if err != nil {
		log.Printf("repository - FindByID(product): %v", err)
		return nil, errors.New("failed to find product by id")
	}
	p := dao.toDomain()
	return &p, nil
}

func (r *ProductRepository) Create(codigo string, idCategoria int, descripcion string, stock int, precioCompra, precioVenta float64) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO productos (codigo, id_categoria, descripcion, stock, precio_compra, precio_venta) VALUES (?, ?, ?, ?, ?, ?)",
		codigo, idCategoria, descripcion, stock, precioCompra, precioVenta,
	)
	if err != nil {
		log.Printf("repository - Create(product): %v", err)
		return 0, errors.New("failed to create product")
	}
	return result.LastInsertId()
}

func (r *ProductRepository) Update(id int, idCategoria int, descripcion string, precioCompra, precioVenta float64) error {
	result, err := r.db.Exec(
		"UPDATE productos SET id_categoria = ?, descripcion = ?, precio_compra = ?, precio_venta = ? WHERE id = ?",
		idCategoria, descripcion, precioCompra, precioVenta, id,
	)
	if err != nil {
		log.Printf("repository - Update(product): %v", err)
		return errors.New("failed to update product")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("repository - Update(product) rows affected: %v", err)
		return errors.New("failed to confirm product update")
	}
	if rows == 0 {
		return apperror.NotFoundError{Message: "product not found"}
	}
	return nil
}

func (r *ProductRepository) UpdateImage(id int, imagen string) error {
	_, err := r.db.Exec("UPDATE productos SET imagen = ? WHERE id = ?", imagen, id)
	if err != nil {
		log.Printf("repository - UpdateImage: %v", err)
		return errors.New("failed to update product image")
	}
	return nil
}

func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec("UPDATE productos SET estado = 0 WHERE id = ?", id)
	if err != nil {
		log.Printf("repository - Delete(product): %v", err)
		return errors.New("failed to delete product")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("repository - Delete(product) rows affected: %v", err)
		return errors.New("failed to confirm product deletion")
	}
	if rows == 0 {
		return apperror.NotFoundError{Message: "product not found"}
	}
	return nil
}

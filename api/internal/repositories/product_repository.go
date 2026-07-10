package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll(categoryID int) ([]models.Product, error) {
	var products []models.Product
	var err error
	if categoryID > 0 {
		err = r.db.Select(&products, "SELECT * FROM productos WHERE id_categoria = ?", categoryID)
	} else {
		err = r.db.Select(&products, "SELECT * FROM productos")
	}
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) FindByID(id int) (*models.Product, error) {
	var product models.Product
	err := r.db.Get(&product, "SELECT * FROM productos WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Create(idCategoria int, descripcion string, stock int, precioCompra, precioVenta float64) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO productos (id_categoria, descripcion, stock, precio_compra, precio_venta) VALUES (?, ?, ?, ?, ?)",
		idCategoria, descripcion, stock, precioCompra, precioVenta,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ProductRepository) Update(id int, idCategoria int, descripcion string, precioCompra, precioVenta float64) error {
	result, err := r.db.Exec(
		"UPDATE productos SET id_categoria = ?, descripcion = ?, precio_compra = ?, precio_venta = ? WHERE id = ?",
		idCategoria, descripcion, precioCompra, precioVenta, id,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *ProductRepository) UpdateImage(id int, imagen string) error {
	_, err := r.db.Exec("UPDATE productos SET imagen = ? WHERE id = ?", imagen, id)
	return err
}

func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec("UPDATE productos SET estado = 0 WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

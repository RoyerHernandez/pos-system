package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Select(&categories, "SELECT * FROM categorias")
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) FindByID(id int) (*models.Category, error) {
	var category models.Category
	err := r.db.Get(&category, "SELECT * FROM categorias WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(nombre, descripcion string) (int64, error) {
	result, err := r.db.Exec("INSERT INTO categorias (nombre, descripcion) VALUES (?, ?)", nombre, descripcion)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *CategoryRepository) Update(id int, nombre, descripcion string) error {
	result, err := r.db.Exec("UPDATE categorias SET nombre = ?, descripcion = ? WHERE id = ?", nombre, descripcion, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	result, err := r.db.Exec("UPDATE categorias SET estado = 0 WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

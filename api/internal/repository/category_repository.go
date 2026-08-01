package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) (*CategoryRepository, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &CategoryRepository{db: db}, nil
}

func (r *CategoryRepository) FindAll() ([]domain.Category, error) {
	var daos []categoryDAO
	err := r.db.Select(&daos, "SELECT * FROM categorias")
	if err != nil {
		log.Printf("repository - FindAll(categories): %v", err)
		return nil, errors.New("failed to find categories")
	}
	categories := make([]domain.Category, len(daos))
	for i := range daos {
		categories[i] = daos[i].toDomain()
	}
	return categories, nil
}

func (r *CategoryRepository) FindByID(id int) (*domain.Category, error) {
	var dao categoryDAO
	err := r.db.Get(&dao, "SELECT * FROM categorias WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFoundError{Message: "category not found"}
	}
	if err != nil {
		log.Printf("repository - FindByID(category): %v", err)
		return nil, errors.New("failed to find category by id")
	}
	c := dao.toDomain()
	return &c, nil
}

func (r *CategoryRepository) Create(nombre, descripcion string) (int64, error) {
	result, err := r.db.Exec("INSERT INTO categorias (nombre, descripcion) VALUES (?, ?)", nombre, descripcion)
	if err != nil {
		log.Printf("repository - Create(category): %v", err)
		return 0, errors.New("failed to create category")
	}
	return result.LastInsertId()
}

func (r *CategoryRepository) Update(id int, nombre, descripcion string) error {
	result, err := r.db.Exec("UPDATE categorias SET nombre = ?, descripcion = ? WHERE id = ?", nombre, descripcion, id)
	if err != nil {
		log.Printf("repository - Update(category): %v", err)
		return errors.New("failed to update category")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("repository - Update(category) rows affected: %v", err)
		return errors.New("failed to confirm category update")
	}
	if rows == 0 {
		return apperror.NotFoundError{Message: "category not found"}
	}
	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	result, err := r.db.Exec("UPDATE categorias SET estado = 0 WHERE id = ?", id)
	if err != nil {
		log.Printf("repository - Delete(category): %v", err)
		return errors.New("failed to delete category")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("repository - Delete(category) rows affected: %v", err)
		return errors.New("failed to confirm category deletion")
	}
	if rows == 0 {
		return apperror.NotFoundError{Message: "category not found"}
	}
	return nil
}

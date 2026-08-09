package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

type TableRepository struct {
	db *sqlx.DB
}

func NewTableRepository(db *sqlx.DB) (*TableRepository, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &TableRepository{db: db}, nil
}

func (r *TableRepository) FindAll() ([]domain.Table, error) {
	var daos []tableDAO
	err := r.db.Select(&daos, "SELECT * FROM mesas ORDER BY numero ASC")
	if err != nil {
		log.Printf("repository - FindAll(tables): %v", err)
		return nil, errors.New("failed to find tables")
	}
	tables := make([]domain.Table, len(daos))
	for i := range daos {
		tables[i] = daos[i].toDomain()
	}
	return tables, nil
}

func (r *TableRepository) FindByID(id int) (*domain.Table, error) {
	var dao tableDAO
	err := r.db.Get(&dao, "SELECT * FROM mesas WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFoundError{Message: "table not found"}
	}
	if err != nil {
		log.Printf("repository - FindByID(table): %v", err)
		return nil, errors.New("failed to find table by id")
	}
	t := dao.toDomain()
	return &t, nil
}

func (r *TableRepository) Create(t *domain.Table) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO mesas (numero, nombre, capacidad) VALUES (?, ?, ?)",
		t.Numero, t.Nombre, t.Capacidad,
	)
	if err != nil {
		log.Printf("repository - Create(table): %v", err)
		return 0, errors.New("failed to create table")
	}
	return result.LastInsertId()
}

func (r *TableRepository) Update(t *domain.Table) error {
	result, err := r.db.Exec(
		"UPDATE mesas SET numero = ?, nombre = ?, capacidad = ? WHERE id = ?",
		t.Numero, t.Nombre, t.Capacidad, t.ID,
	)
	if err != nil {
		log.Printf("repository - Update(table): %v", err)
		return errors.New("failed to update table")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("repository - Update(table) rows affected: %v", err)
		return errors.New("failed to confirm table update")
	}
	if rows == 0 {
		return apperror.NotFoundError{Message: "table not found"}
	}
	return nil
}

func (r *TableRepository) Delete(id int) error {
	// Only allow deletion if the table is free
	var estado string
	err := r.db.Get(&estado, "SELECT estado FROM mesas WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return apperror.NotFoundError{Message: "table not found"}
	}
	if err != nil {
		log.Printf("repository - Delete(table) check state: %v", err)
		return errors.New("failed to check table state")
	}
	if estado != "libre" {
		return apperror.ValidationError{Message: "cannot delete a table that is not free"}
	}

	result, err := r.db.Exec("DELETE FROM mesas WHERE id = ?", id)
	if err != nil {
		log.Printf("repository - Delete(table): %v", err)
		return errors.New("failed to delete table")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("repository - Delete(table) rows affected: %v", err)
		return errors.New("failed to confirm table deletion")
	}
	if rows == 0 {
		return apperror.NotFoundError{Message: "table not found"}
	}
	return nil
}

func (r *TableRepository) ExistsByNumero(numero int, excludeID int) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM mesas WHERE numero = ? AND id != ?", numero, excludeID)
	if err != nil {
		log.Printf("repository - ExistsByNumero(table): %v", err)
		return false, errors.New("failed to check table number")
	}
	return count > 0, nil
}

// FindByIDForUpdate locks the row with SELECT ... FOR UPDATE within a transaction.
func (r *TableRepository) FindByIDForUpdate(tx *sqlx.Tx, id int) (*domain.Table, error) {
	var dao tableDAO
	err := tx.Get(&dao, "SELECT * FROM mesas WHERE id = ? FOR UPDATE", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFoundError{Message: "table not found"}
	}
	if err != nil {
		log.Printf("repository - FindByIDForUpdate(table): %v", err)
		return nil, errors.New("failed to find table for update")
	}
	t := dao.toDomain()
	return &t, nil
}

// UpdateState changes the estado, id_venta_activa and id_mesero of a table within a transaction.
func (r *TableRepository) UpdateState(tx *sqlx.Tx, id int, estado string, idVenta *int, idMesero *int) error {
	_, err := tx.Exec(
		"UPDATE mesas SET estado = ?, id_venta_activa = ?, id_mesero = ?, fecha_actualizacion = NOW() WHERE id = ?",
		estado, idVenta, idMesero, id,
	)
	if err != nil {
		log.Printf("repository - UpdateState(table): %v", err)
		return errors.New("failed to update table state")
	}
	return nil
}

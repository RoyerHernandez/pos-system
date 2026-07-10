package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/jmoiron/sqlx"
)

type ClientRepository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) FindAll() ([]models.Client, error) {
	var clients []models.Client
	err := r.db.Select(&clients, "SELECT * FROM clientes")
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *ClientRepository) FindByID(id int) (*models.Client, error) {
	var client models.Client
	err := r.db.Get(&client, "SELECT * FROM clientes WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) Create(nombre string, documento, email, telefono, direccion, fechaNacimiento *string) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO clientes (nombre, documento, email, telefono, direccion, fecha_nacimiento) VALUES (?, ?, ?, ?, ?, ?)",
		nombre, documento, email, telefono, direccion, fechaNacimiento,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ClientRepository) Update(id int, nombre string, documento, email, telefono, direccion, fechaNacimiento *string) error {
	result, err := r.db.Exec(
		"UPDATE clientes SET nombre = ?, documento = ?, email = ?, telefono = ?, direccion = ?, fecha_nacimiento = ? WHERE id = ?",
		nombre, documento, email, telefono, direccion, fechaNacimiento, id,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("client not found")
	}
	return nil
}

func (r *ClientRepository) Delete(id int) error {
	result, err := r.db.Exec("UPDATE clientes SET estado = 0 WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("client not found")
	}
	return nil
}

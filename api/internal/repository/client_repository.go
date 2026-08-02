package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

type ClientRepository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) (*ClientRepository, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &ClientRepository{db: db}, nil
}

func (r *ClientRepository) FindAll() ([]domain.Client, error) {
	var daos []clientDAO
	err := r.db.Select(&daos, "SELECT * FROM clientes")
	if err != nil {
		log.Printf("repository - FindAll(clients): %v", err)
		return nil, errors.New("failed to find clients")
	}
	clients := make([]domain.Client, len(daos))
	for i := range daos {
		clients[i] = daos[i].toDomain()
	}
	return clients, nil
}

func (r *ClientRepository) FindByID(id int) (*domain.Client, error) {
	var dao clientDAO
	err := r.db.Get(&dao, "SELECT * FROM clientes WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFoundError{Message: "client not found"}
	}
	if err != nil {
		log.Printf("repository - FindByID(client): %v", err)
		return nil, errors.New("failed to find client by id")
	}
	c := dao.toDomain()
	return &c, nil
}

func (r *ClientRepository) Create(nombre string, documento, email, telefono, direccion, fechaNacimiento *string) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO clientes (nombre, documento, email, telefono, direccion, fecha_nacimiento) VALUES (?, ?, ?, ?, ?, ?)",
		nombre, documento, email, telefono, direccion, fechaNacimiento,
	)
	if err != nil {
		log.Printf("repository - Create(client): %v", err)
		return 0, errors.New("failed to create client")
	}
	return result.LastInsertId()
}

func (r *ClientRepository) Update(id int, nombre string, documento, email, telefono, direccion, fechaNacimiento *string) error {
	result, err := r.db.Exec(
		"UPDATE clientes SET nombre = ?, documento = ?, email = ?, telefono = ?, direccion = ?, fecha_nacimiento = ? WHERE id = ?",
		nombre, documento, email, telefono, direccion, fechaNacimiento, id,
	)
	if err != nil {
		log.Printf("repository - Update(client): %v", err)
		return errors.New("failed to update client")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("repository - Update(client) rows affected: %v", err)
		return errors.New("failed to confirm client update")
	}
	if rows == 0 {
		return apperror.NotFoundError{Message: "client not found"}
	}
	return nil
}

func (r *ClientRepository) Delete(id int) error {
	result, err := r.db.Exec("UPDATE clientes SET estado = 0 WHERE id = ?", id)
	if err != nil {
		log.Printf("repository - Delete(client): %v", err)
		return errors.New("failed to delete client")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("repository - Delete(client) rows affected: %v", err)
		return errors.New("failed to confirm client deletion")
	}
	if rows == 0 {
		return apperror.NotFoundError{Message: "client not found"}
	}
	return nil
}

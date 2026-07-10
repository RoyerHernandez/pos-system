package repositories

import (
	"database/sql"
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM usuarios WHERE usuario = ? AND estado = 1", username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Select(&users, "SELECT * FROM usuarios WHERE estado = 1")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM usuarios WHERE id = ? AND estado = 1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO usuarios (usuario, password, nombre, perfil, foto, estado) VALUES (?, ?, ?, ?, ?, ?)",
		user.Usuario, user.Password, user.Nombre, user.Perfil, user.Foto, user.Estado,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *UserRepository) Update(user *models.User) error {
	if user.Password != "" {
		_, err := r.db.Exec(
			"UPDATE usuarios SET usuario = ?, password = ?, nombre = ?, perfil = ?, foto = ? WHERE id = ?",
			user.Usuario, user.Password, user.Nombre, user.Perfil, user.Foto, user.ID,
		)
		return err
	}
	_, err := r.db.Exec(
		"UPDATE usuarios SET usuario = ?, nombre = ?, perfil = ?, foto = ? WHERE id = ?",
		user.Usuario, user.Nombre, user.Perfil, user.Foto, user.ID,
	)
	return err
}

func (r *UserRepository) SoftDelete(id int) error {
	_, err := r.db.Exec("UPDATE usuarios SET estado = 0 WHERE id = ?", id)
	return err
}

func (r *UserRepository) UpdateLastLogin(id int) error {
	_, err := r.db.Exec("UPDATE usuarios SET ultimo_login = NOW() WHERE id = ?", id)
	return err
}

func (r *UserRepository) ExistsByUsername(username string, excludeID int) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM usuarios WHERE usuario = ? AND id != ? AND estado = 1", username, excludeID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) (*UserRepository, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	var dao userDAO
	err := r.db.Get(&dao, "SELECT * FROM usuarios WHERE usuario = ? AND estado = 1", username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFoundError{Message: "user not found"}
	}
	if err != nil {
		log.Printf("repository - FindByUsername: %v", err)
		return nil, errors.New("failed to find user by username")
	}
	u := dao.toDomain()
	return &u, nil
}

func (r *UserRepository) FindAll() ([]domain.User, error) {
	var daos []userDAO
	err := r.db.Select(&daos, "SELECT * FROM usuarios WHERE estado = 1")
	if err != nil {
		log.Printf("repository - FindAll(users): %v", err)
		return nil, errors.New("failed to find users")
	}
	users := make([]domain.User, len(daos))
	for i := range daos {
		users[i] = daos[i].toDomain()
	}
	return users, nil
}

func (r *UserRepository) FindByID(id int) (*domain.User, error) {
	var dao userDAO
	err := r.db.Get(&dao, "SELECT * FROM usuarios WHERE id = ? AND estado = 1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFoundError{Message: "user not found"}
	}
	if err != nil {
		log.Printf("repository - FindByID(user): %v", err)
		return nil, errors.New("failed to find user by id")
	}
	u := dao.toDomain()
	return &u, nil
}

func (r *UserRepository) Create(user *domain.User) (int64, error) {
	dao := toUserDAO(*user)
	result, err := r.db.Exec(
		"INSERT INTO usuarios (usuario, password, nombre, perfil, foto, estado) VALUES (?, ?, ?, ?, ?, ?)",
		dao.Usuario, dao.Password, dao.Nombre, dao.Perfil, dao.Foto, dao.Estado,
	)
	if err != nil {
		log.Printf("repository - Create(user): %v", err)
		return 0, errors.New("failed to create user")
	}
	return result.LastInsertId()
}

func (r *UserRepository) Update(user *domain.User) error {
	dao := toUserDAO(*user)
	if dao.Password != "" {
		_, err := r.db.Exec(
			"UPDATE usuarios SET usuario = ?, password = ?, nombre = ?, perfil = ?, foto = ? WHERE id = ?",
			dao.Usuario, dao.Password, dao.Nombre, dao.Perfil, dao.Foto, dao.ID,
		)
		if err != nil {
			log.Printf("repository - Update(user): %v", err)
			return errors.New("failed to update user")
		}
		return nil
	}
	_, err := r.db.Exec(
		"UPDATE usuarios SET usuario = ?, nombre = ?, perfil = ?, foto = ? WHERE id = ?",
		dao.Usuario, dao.Nombre, dao.Perfil, dao.Foto, dao.ID,
	)
	if err != nil {
		log.Printf("repository - Update(user): %v", err)
		return errors.New("failed to update user")
	}
	return nil
}

func (r *UserRepository) SoftDelete(id int) error {
	_, err := r.db.Exec("UPDATE usuarios SET estado = 0 WHERE id = ?", id)
	if err != nil {
		log.Printf("repository - SoftDelete(user): %v", err)
		return errors.New("failed to soft delete user")
	}
	return nil
}

func (r *UserRepository) UpdateLastLogin(id int) error {
	_, err := r.db.Exec("UPDATE usuarios SET ultimo_login = NOW() WHERE id = ?", id)
	if err != nil {
		log.Printf("repository - UpdateLastLogin: %v", err)
		return errors.New("failed to update last login")
	}
	return nil
}

func (r *UserRepository) ExistsByUsername(username string, excludeID int) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM usuarios WHERE usuario = ? AND id != ? AND estado = 1", username, excludeID)
	if err != nil {
		log.Printf("repository - ExistsByUsername: %v", err)
		return false, errors.New("failed to check username existence")
	}
	return count > 0, nil
}

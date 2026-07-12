package service

import (
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var validProfiles = map[string]bool{
	"Administrador": true,
	"Especial":      true,
	"Vendedor":      true,
}

type userRepository interface {
	FindAll() ([]domain.User, error)
	FindByID(id int) (*domain.User, error)
	Create(user *domain.User) (int64, error)
	Update(user *domain.User) error
	SoftDelete(id int) error
	ExistsByUsername(username string, excludeID int) (bool, error)
}

type UserService struct {
	users userRepository
}

func NewUserService(users userRepository) (*UserService, error) {
	if users == nil {
		return nil, errors.New("userRepository must not be nil")
	}
	return &UserService{users: users}, nil
}

func (s *UserService) GetAll() ([]domain.User, error) {
	return s.users.FindAll()
}

func (s *UserService) GetByID(id int) (*domain.User, error) {
	return s.users.FindByID(id)
}

func (s *UserService) Create(user *domain.User, password string) (int64, error) {
	if user.Usuario == "" {
		return 0, errors.New("usuario is required")
	}
	if !validProfiles[user.Perfil] {
		return 0, errors.New("perfil must be Administrador, Especial, or Vendedor")
	}
	if password == "" {
		return 0, errors.New("password is required")
	}
	if len(password) < 6 {
		return 0, errors.New("password must be at least 6 characters")
	}

	exists, err := s.users.ExistsByUsername(user.Usuario, 0)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("usuario already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("failed to hash password")
	}
	user.Password = string(hashed)
	user.Estado = 1

	return s.users.Create(user)
}

func (s *UserService) Update(user *domain.User, password string) error {
	if user.Usuario == "" {
		return errors.New("usuario is required")
	}
	if !validProfiles[user.Perfil] {
		return errors.New("perfil must be Administrador, Especial, or Vendedor")
	}

	exists, err := s.users.ExistsByUsername(user.Usuario, user.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("usuario already exists")
	}

	if password != "" {
		if len(password) < 6 {
			return errors.New("password must be at least 6 characters")
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}
		user.Password = string(hashed)
	}

	return s.users.Update(user)
}

func (s *UserService) Delete(id int) error {
	return s.users.SoftDelete(id)
}

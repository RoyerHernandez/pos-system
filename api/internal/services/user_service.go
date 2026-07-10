package services

import (
	"errors"

	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var validProfiles = map[string]bool{
	"Administrador": true,
	"Especial":      true,
	"Vendedor":      true,
}

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) Create(user *models.User, password string) (int64, error) {
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

	exists, err := s.userRepo.ExistsByUsername(user.Usuario, 0)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("usuario already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(hashed)
	user.Estado = 1

	return s.userRepo.Create(user)
}

func (s *UserService) Update(user *models.User, password string) error {
	if user.Usuario == "" {
		return errors.New("usuario is required")
	}
	if !validProfiles[user.Perfil] {
		return errors.New("perfil must be Administrador, Especial, or Vendedor")
	}

	exists, err := s.userRepo.ExistsByUsername(user.Usuario, user.ID)
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
			return err
		}
		user.Password = string(hashed)
	}

	return s.userRepo.Update(user)
}

func (s *UserService) Delete(id int) error {
	return s.userRepo.SoftDelete(id)
}

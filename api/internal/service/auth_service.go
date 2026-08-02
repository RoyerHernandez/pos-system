package service

import (
	"errors"
	"time"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/pkg/infrastructure/auth"
	"golang.org/x/crypto/bcrypt"
)

type authUserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	FindByID(id int) (*domain.User, error)
	UpdateLastLogin(id int) error
}

type AuthService struct {
	users     authUserRepository
	jwtSecret string
}

func NewAuthService(users authUserRepository, jwtSecret string) (*AuthService, error) {
	if users == nil {
		return nil, errors.New("userRepository must not be nil")
	}
	if jwtSecret == "" {
		return nil, errors.New("jwtSecret must not be empty")
	}
	return &AuthService{users: users, jwtSecret: jwtSecret}, nil
}

type AuthResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	User         domain.UserResponse `json:"user"`
}

func (s *AuthService) Login(username, password string) (*AuthResponse, error) {
	user, err := s.users.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Update last login timestamp (non-blocking: don't fail login if this fails)
	_ = s.users.UpdateLastLogin(user.ID)

	accessToken, err := auth.GenerateToken(user.ID, user.Perfil, s.jwtSecret, 24*time.Hour, "access")
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := auth.GenerateToken(user.ID, user.Perfil, s.jwtSecret, 7*24*time.Hour, "refresh")
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	claims, err := auth.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.users.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user is deactivated or not found")
	}
	if user == nil || user.Estado != 1 {
		return nil, errors.New("user is deactivated or not found")
	}

	accessToken, err := auth.GenerateToken(claims.UserID, claims.Role, s.jwtSecret, 24*time.Hour, "access")
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	newRefreshToken, err := auth.GenerateToken(claims.UserID, claims.Role, s.jwtSecret, 7*24*time.Hour, "refresh")
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user.ToResponse(),
	}, nil
}

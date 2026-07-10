package services

import (
	"errors"
	"time"

	"github.com/RoyerHernandez/pos-system/api/internal/auth"
	"github.com/RoyerHernandez/pos-system/api/internal/models"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

type AuthResponse struct {
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	User         models.UserResponse `json:"user"`
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *AuthService) Login(username, password string) (*AuthResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Update last login timestamp (non-blocking: don't fail login if this fails)
	_ = s.userRepo.UpdateLastLogin(user.ID)

	accessToken, err := auth.GenerateToken(user.ID, user.Perfil, s.jwtSecret, 24*time.Hour, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateToken(user.ID, user.Perfil, s.jwtSecret, 7*24*time.Hour, "refresh")
	if err != nil {
		return nil, err
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

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Estado != 1 {
		return nil, errors.New("user is deactivated or not found")
	}

	accessToken, err := auth.GenerateToken(claims.UserID, claims.Role, s.jwtSecret, 24*time.Hour, "access")
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := auth.GenerateToken(claims.UserID, claims.Role, s.jwtSecret, 7*24*time.Hour, "refresh")
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user.ToResponse(),
	}, nil
}

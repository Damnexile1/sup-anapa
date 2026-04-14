package services

import (
	"context"
	"errors"
	"sup-anapa/internal/models"
	"sup-anapa/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type AuthService struct {
	adminRepo *repository.AdminRepository
}

func NewAuthService(adminRepo *repository.AdminRepository) *AuthService {
	return &AuthService{
		adminRepo: adminRepo,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, username, password string) (*models.Admin, error) {
	admin, err := s.adminRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return admin, nil
}

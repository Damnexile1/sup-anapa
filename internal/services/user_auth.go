package services

import (
	"context"
	"errors"
	"sup-anapa/internal/models"
	"sup-anapa/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserAuthService struct {
	repo *repository.UserRepository
}

func NewUserAuthService(repo *repository.UserRepository) *UserAuthService {
	return &UserAuthService{repo: repo}
}

func (s *UserAuthService) Register(ctx context.Context, username, password, phone string) (*models.User, error) {
	if username == "" || password == "" || phone == "" {
		return nil, errors.New("all fields are required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{
		Username:     username,
		PasswordHash: string(hash),
		Phone:        phone,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserAuthService) Login(ctx context.Context, username, password string) (*models.User, error) {
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}

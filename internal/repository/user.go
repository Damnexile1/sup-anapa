package repository

import (
	"context"
	"sup-anapa/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, password_hash, phone) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, user.Username, user.PasswordHash, user.Phone).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, password_hash, phone, created_at FROM users WHERE username = $1`
	u := &models.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Phone, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, username, password_hash, phone, created_at FROM users WHERE id = $1`
	u := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Phone, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

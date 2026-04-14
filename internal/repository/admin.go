package repository

import (
	"context"
	"sup-anapa/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) Create(ctx context.Context, admin *models.Admin) error {
	query := `
		INSERT INTO admins (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query, admin.Username, admin.PasswordHash).Scan(&admin.ID, &admin.CreatedAt)
}

func (r *AdminRepository) GetByUsername(ctx context.Context, username string) (*models.Admin, error) {
	query := `SELECT id, username, password_hash, created_at FROM admins WHERE username = $1`

	admin := &models.Admin{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (r *AdminRepository) GetByID(ctx context.Context, id int) (*models.Admin, error) {
	query := `SELECT id, username, password_hash, created_at FROM admins WHERE id = $1`

	admin := &models.Admin{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

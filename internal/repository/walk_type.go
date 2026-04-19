package repository

import (
	"context"
	"sup-anapa/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WalkTypeRepository struct {
	db *pgxpool.Pool
}

func NewWalkTypeRepository(db *pgxpool.Pool) *WalkTypeRepository {
	return &WalkTypeRepository{db: db}
}

func (r *WalkTypeRepository) Create(ctx context.Context, wt *models.WalkType) error {
	query := `
		INSERT INTO walk_types (instructor_id, name, price, max_people)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, wt.InstructorID, wt.Name, wt.Price, wt.MaxPeople).
		Scan(&wt.ID, &wt.CreatedAt, &wt.UpdatedAt)
}

func (r *WalkTypeRepository) GetByID(ctx context.Context, id int) (*models.WalkType, error) {
	query := `SELECT id, instructor_id, name, price, max_people, created_at, updated_at FROM walk_types WHERE id = $1`
	wt := &models.WalkType{}
	err := r.db.QueryRow(ctx, query, id).Scan(&wt.ID, &wt.InstructorID, &wt.Name, &wt.Price, &wt.MaxPeople, &wt.CreatedAt, &wt.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return wt, nil
}

func (r *WalkTypeRepository) GetByInstructorID(ctx context.Context, instructorID int) ([]*models.WalkType, error) {
	query := `SELECT id, instructor_id, name, price, max_people, created_at, updated_at FROM walk_types WHERE instructor_id = $1 ORDER BY name`
	rows, err := r.db.Query(ctx, query, instructorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var walkTypes []*models.WalkType
	for rows.Next() {
		wt := &models.WalkType{}
		if err := rows.Scan(&wt.ID, &wt.InstructorID, &wt.Name, &wt.Price, &wt.MaxPeople, &wt.CreatedAt, &wt.UpdatedAt); err != nil {
			return nil, err
		}
		walkTypes = append(walkTypes, wt)
	}
	return walkTypes, nil
}

func (r *WalkTypeRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM walk_types WHERE id = $1`, id)
	return err
}

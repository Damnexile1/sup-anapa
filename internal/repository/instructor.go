package repository

import (
	"context"
	"sup-anapa/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InstructorRepository struct {
	db *pgxpool.Pool
}

func NewInstructorRepository(db *pgxpool.Pool) *InstructorRepository {
	return &InstructorRepository{db: db}
}

func (r *InstructorRepository) Create(ctx context.Context, instructor *models.Instructor) error {
	query := `
		INSERT INTO instructors (name, photo, description, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		ctx,
		query,
		instructor.Name,
		instructor.Photo,
		instructor.Description,
		instructor.Phone,
	).Scan(&instructor.ID, &instructor.CreatedAt, &instructor.UpdatedAt)
}

func (r *InstructorRepository) GetByID(ctx context.Context, id int) (*models.Instructor, error) {
	query := `SELECT id, name, photo, description, phone, created_at, updated_at FROM instructors WHERE id = $1`

	instructor := &models.Instructor{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&instructor.ID,
		&instructor.Name,
		&instructor.Photo,
		&instructor.Description,
		&instructor.Phone,
		&instructor.CreatedAt,
		&instructor.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	walkTypes, err := r.getWalkTypes(ctx, instructor.ID)
	if err == nil {
		instructor.WalkTypes = walkTypes
	}

	return instructor, nil
}

func (r *InstructorRepository) GetAll(ctx context.Context) ([]*models.Instructor, error) {
	query := `SELECT id, name, photo, description, phone, created_at, updated_at FROM instructors ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instructors []*models.Instructor
	for rows.Next() {
		instructor := &models.Instructor{}
		err := rows.Scan(
			&instructor.ID,
			&instructor.Name,
			&instructor.Photo,
			&instructor.Description,
			&instructor.Phone,
			&instructor.CreatedAt,
			&instructor.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		walkTypes, err := r.getWalkTypes(ctx, instructor.ID)
		if err == nil {
			instructor.WalkTypes = walkTypes
		}
		instructors = append(instructors, instructor)
	}
	return instructors, nil
}

func (r *InstructorRepository) Update(ctx context.Context, instructor *models.Instructor) error {
	query := `
		UPDATE instructors 
		SET name = $1, photo = $2, description = $3, phone = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query, instructor.Name, instructor.Photo, instructor.Description, instructor.Phone, instructor.ID)
	return err
}

func (r *InstructorRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM instructors WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *InstructorRepository) getWalkTypes(ctx context.Context, instructorID int) ([]*models.WalkType, error) {
	query := `SELECT id, instructor_id, name, price, max_people, created_at, updated_at FROM walk_types WHERE instructor_id = $1 ORDER BY name`
	rows, err := r.db.Query(ctx, query, instructorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	walkTypes := make([]*models.WalkType, 0)
	for rows.Next() {
		wt := &models.WalkType{}
		if err := rows.Scan(&wt.ID, &wt.InstructorID, &wt.Name, &wt.Price, &wt.MaxPeople, &wt.CreatedAt, &wt.UpdatedAt); err != nil {
			return nil, err
		}
		walkTypes = append(walkTypes, wt)
	}

	return walkTypes, nil
}

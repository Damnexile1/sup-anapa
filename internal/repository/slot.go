package repository

import (
	"context"
	"sup-anapa/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SlotRepository struct {
	db *pgxpool.Pool
}

func NewSlotRepository(db *pgxpool.Pool) *SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) Create(ctx context.Context, slot *models.Slot) error {
	query := `
		INSERT INTO slots (date, start_time, end_time, price, max_people, instructor_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'available')
		RETURNING id, created_at, updated_at
	`
	slot.Status = "available"
	return r.db.QueryRow(
		ctx,
		query,
		slot.Date,
		slot.StartTime,
		slot.EndTime,
		slot.Price,
		slot.MaxPeople,
		slot.InstructorID,
	).Scan(&slot.ID, &slot.CreatedAt, &slot.UpdatedAt)
}

func (r *SlotRepository) GetByID(ctx context.Context, id int) (*models.Slot, error) {
	query := `SELECT id, date, start_time, end_time, price, max_people, instructor_id, status, hold_expires_at, created_at, updated_at
			  FROM slots WHERE id = $1`

	slot := &models.Slot{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&slot.ID,
		&slot.Date,
		&slot.StartTime,
		&slot.EndTime,
		&slot.Price,
		&slot.MaxPeople,
		&slot.InstructorID,
		&slot.Status,
		&slot.HoldExpiresAt,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return slot, nil
}

func (r *SlotRepository) GetAll(ctx context.Context) ([]*models.Slot, error) {
	query := `SELECT id, date, start_time, end_time, price, max_people, instructor_id, status, hold_expires_at, created_at, updated_at
			  FROM slots ORDER BY date, start_time`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*models.Slot
	for rows.Next() {
		slot := &models.Slot{}
		err := rows.Scan(
			&slot.ID,
			&slot.Date,
			&slot.StartTime,
			&slot.EndTime,
			&slot.Price,
			&slot.MaxPeople,
			&slot.InstructorID,
			&slot.Status,
			&slot.HoldExpiresAt,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *SlotRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Slot, error) {
	query := `SELECT id, date, start_time, end_time, price, max_people, instructor_id, status, hold_expires_at, created_at, updated_at
			  FROM slots WHERE date BETWEEN $1 AND $2 ORDER BY date, start_time`

	rows, err := r.db.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*models.Slot
	for rows.Next() {
		slot := &models.Slot{}
		err := rows.Scan(
			&slot.ID,
			&slot.Date,
			&slot.StartTime,
			&slot.EndTime,
			&slot.Price,
			&slot.MaxPeople,
			&slot.InstructorID,
			&slot.Status,
			&slot.HoldExpiresAt,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *SlotRepository) GetAvailableSlots(ctx context.Context, start, end time.Time) ([]*models.Slot, error) {
	query := `
		SELECT s.id, s.date, s.start_time, s.end_time, s.price, s.max_people, s.instructor_id, s.status, s.hold_expires_at, s.created_at, s.updated_at
		FROM slots s
		WHERE s.date BETWEEN $1 AND $2
		AND s.date >= CURRENT_DATE
		AND s.status = 'available'
		AND (
			SELECT COUNT(*) FROM bookings b
			WHERE b.slot_id = s.id AND b.status != 'cancelled'
		) < s.max_people
		ORDER BY s.date, s.start_time
	`

	rows, err := r.db.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*models.Slot
	for rows.Next() {
		slot := &models.Slot{}
		err := rows.Scan(
			&slot.ID,
			&slot.Date,
			&slot.StartTime,
			&slot.EndTime,
			&slot.Price,
			&slot.MaxPeople,
			&slot.InstructorID,
			&slot.Status,
			&slot.HoldExpiresAt,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *SlotRepository) Update(ctx context.Context, slot *models.Slot) error {
	query := `
		UPDATE slots
		SET date = $1, start_time = $2, end_time = $3, price = $4, max_people = $5, instructor_id = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`
	_, err := r.db.Exec(ctx, query, slot.Date, slot.StartTime, slot.EndTime, slot.Price, slot.MaxPeople, slot.InstructorID, slot.ID)
	return err
}

func (r *SlotRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM slots WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *SlotRepository) SetPending(ctx context.Context, slotID int, holdExpiresAt time.Time) error {
	query := `UPDATE slots SET status = 'pending', hold_expires_at = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(ctx, query, holdExpiresAt, slotID)
	return err
}

func (r *SlotRepository) SetAvailable(ctx context.Context, slotID int) error {
	query := `UPDATE slots SET status = 'available', hold_expires_at = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(ctx, query, slotID)
	return err
}

func (r *SlotRepository) SetBooked(ctx context.Context, slotID int) error {
	query := `UPDATE slots SET status = 'booked', hold_expires_at = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(ctx, query, slotID)
	return err
}

func (r *SlotRepository) ExpireHolds(ctx context.Context) (int, error) {
	query := `UPDATE slots SET status = 'available', hold_expires_at = NULL, updated_at = CURRENT_TIMESTAMP
			  WHERE status = 'pending' AND hold_expires_at < CURRENT_TIMESTAMP`
	cmd, err := r.db.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return int(cmd.RowsAffected()), nil
}

func (r *SlotRepository) GetByIDWithLock(ctx context.Context, id int) (*models.Slot, error) {
	query := `SELECT id, date, start_time, end_time, price, max_people, instructor_id, status, hold_expires_at, created_at, updated_at
			  FROM slots WHERE id = $1 FOR UPDATE`

	slot := &models.Slot{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&slot.ID,
		&slot.Date,
		&slot.StartTime,
		&slot.EndTime,
		&slot.Price,
		&slot.MaxPeople,
		&slot.InstructorID,
		&slot.Status,
		&slot.HoldExpiresAt,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return slot, nil
}

package repository

import (
	"context"
	"strconv"
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

const slotSelectColumns = `s.id, s.date, s.start_time, s.end_time, s.price, s.max_people, s.instructor_id, s.walk_type_id, COALESCE(wt.name, ''), s.status, s.hold_expires_at, s.created_at, s.updated_at`

func scanSlot(scanner interface {
	Scan(dest ...interface{}) error
}, slot *models.Slot) error {
	return scanner.Scan(
		&slot.ID,
		&slot.Date,
		&slot.StartTime,
		&slot.EndTime,
		&slot.Price,
		&slot.MaxPeople,
		&slot.InstructorID,
		&slot.WalkTypeID,
		&slot.WalkTypeName,
		&slot.Status,
		&slot.HoldExpiresAt,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	)
}

func (r *SlotRepository) Create(ctx context.Context, slot *models.Slot) error {
	query := `
		INSERT INTO slots (date, start_time, end_time, price, max_people, instructor_id, walk_type_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'available')
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
		slot.WalkTypeID,
	).Scan(&slot.ID, &slot.CreatedAt, &slot.UpdatedAt)
}

func (r *SlotRepository) GetByID(ctx context.Context, id int) (*models.Slot, error) {
	query := `SELECT ` + slotSelectColumns + `
			  FROM slots s
			  LEFT JOIN walk_types wt ON wt.id = s.walk_type_id
			  WHERE s.id = $1`

	slot := &models.Slot{}
	err := scanSlot(r.db.QueryRow(ctx, query, id), slot)
	if err != nil {
		return nil, err
	}
	return slot, nil
}

func (r *SlotRepository) GetAll(ctx context.Context) ([]*models.Slot, error) {
	query := `SELECT ` + slotSelectColumns + `
			  FROM slots s
			  LEFT JOIN walk_types wt ON wt.id = s.walk_type_id
			  ORDER BY s.date, s.start_time`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*models.Slot
	for rows.Next() {
		slot := &models.Slot{}
		if err := scanSlot(rows, slot); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *SlotRepository) GetByFilters(ctx context.Context, instructorID, walkTypeID int) ([]*models.Slot, error) {
	query := `SELECT ` + slotSelectColumns + `
			  FROM slots s
			  LEFT JOIN walk_types wt ON wt.id = s.walk_type_id
			  WHERE s.date >= CURRENT_DATE`
	args := make([]interface{}, 0)
	argPos := 1
	if instructorID > 0 {
		query += ` AND s.instructor_id = $` + strconv.Itoa(argPos)
		args = append(args, instructorID)
		argPos++
	}
	if walkTypeID > 0 {
		query += ` AND s.walk_type_id = $` + strconv.Itoa(argPos)
		args = append(args, walkTypeID)
	}
	query += ` ORDER BY s.date, s.start_time`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*models.Slot
	for rows.Next() {
		slot := &models.Slot{}
		if err := scanSlot(rows, slot); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *SlotRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Slot, error) {
	query := `SELECT ` + slotSelectColumns + `
			  FROM slots s
			  LEFT JOIN walk_types wt ON wt.id = s.walk_type_id
			  WHERE s.date BETWEEN $1 AND $2 ORDER BY s.date, s.start_time`

	rows, err := r.db.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*models.Slot
	for rows.Next() {
		slot := &models.Slot{}
		if err := scanSlot(rows, slot); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *SlotRepository) GetAvailableSlots(ctx context.Context, start, end time.Time) ([]*models.Slot, error) {
	query := `
		SELECT ` + slotSelectColumns + `
		FROM slots s
		LEFT JOIN walk_types wt ON wt.id = s.walk_type_id
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
		if err := scanSlot(rows, slot); err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *SlotRepository) Update(ctx context.Context, slot *models.Slot) error {
	query := `
		UPDATE slots
		SET date = $1, start_time = $2, end_time = $3, price = $4, max_people = $5, instructor_id = $6, walk_type_id = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
	`
	_, err := r.db.Exec(ctx, query, slot.Date, slot.StartTime, slot.EndTime, slot.Price, slot.MaxPeople, slot.InstructorID, slot.WalkTypeID, slot.ID)
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
	query := `SELECT ` + slotSelectColumns + `
			  FROM slots s
			  LEFT JOIN walk_types wt ON wt.id = s.walk_type_id
			  WHERE s.id = $1 FOR UPDATE`

	slot := &models.Slot{}
	err := scanSlot(r.db.QueryRow(ctx, query, id), slot)
	if err != nil {
		return nil, err
	}
	return slot, nil
}

package repository

import (
	"context"
	"sup-anapa/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	query := `
		INSERT INTO bookings (slot_id, client_name, client_phone, client_email, people_count, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		ctx,
		query,
		booking.SlotID,
		booking.ClientName,
		booking.ClientPhone,
		booking.ClientEmail,
		booking.PeopleCount,
		booking.Status,
	).Scan(&booking.ID, &booking.CreatedAt, &booking.UpdatedAt)
}

func (r *BookingRepository) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	query := `SELECT id, slot_id, client_name, client_phone, client_email, people_count, status, created_at, updated_at
			  FROM bookings WHERE id = $1`

	booking := &models.Booking{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&booking.ID,
		&booking.SlotID,
		&booking.ClientName,
		&booking.ClientPhone,
		&booking.ClientEmail,
		&booking.PeopleCount,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (r *BookingRepository) GetByStatus(ctx context.Context, status string) ([]*models.Booking, error) {
	query := `SELECT id, slot_id, client_name, client_phone, client_email, people_count, status, created_at, updated_at
			  FROM bookings WHERE status = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		booking := &models.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.SlotID,
			&booking.ClientName,
			&booking.ClientPhone,
			&booking.ClientEmail,
			&booking.PeopleCount,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (r *BookingRepository) GetAll(ctx context.Context) ([]*models.Booking, error) {
	query := `SELECT id, slot_id, client_name, client_phone, client_email, people_count, status, created_at, updated_at
			  FROM bookings ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		booking := &models.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.SlotID,
			&booking.ClientName,
			&booking.ClientPhone,
			&booking.ClientEmail,
			&booking.PeopleCount,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE bookings SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *BookingRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Booking, error) {
	query := `
		SELECT b.id, b.slot_id, b.client_name, b.client_phone, b.client_email, b.people_count, b.status, b.created_at, b.updated_at
		FROM bookings b
		JOIN slots s ON b.slot_id = s.id
		WHERE s.date BETWEEN $1 AND $2
		ORDER BY s.date, s.start_time
	`

	rows, err := r.db.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		booking := &models.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.SlotID,
			&booking.ClientName,
			&booking.ClientPhone,
			&booking.ClientEmail,
			&booking.PeopleCount,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (r *BookingRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM bookings WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

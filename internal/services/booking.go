package services

import (
	"context"
	"sup-anapa/internal/models"
	"sup-anapa/internal/repository"
	"time"
)

type BookingService struct {
	repo         *repository.BookingRepository
	notification *NotificationService
}

func NewBookingService(repo *repository.BookingRepository, notification *NotificationService) *BookingService {
	return &BookingService{
		repo:         repo,
		notification: notification,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, booking *models.Booking) error {
	// Create booking in database
	if err := s.repo.Create(ctx, booking); err != nil {
		return err
	}

	// Send notification to admin
	slotInfo := "Slot details" // TODO: Get actual slot info
	if err := s.notification.SendBookingNotification(
		booking.ClientName,
		booking.ClientPhone,
		slotInfo,
	); err != nil {
		// Log error but don't fail the booking
		// TODO: Add proper logging
	}

	return nil
}

func (s *BookingService) GetBookingsByStatus(ctx context.Context, status string) ([]*models.Booking, error) {
	return s.repo.GetByStatus(ctx, status)
}

func (s *BookingService) UpdateBookingStatus(ctx context.Context, id int, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *BookingService) GetBookingsByDateRange(ctx context.Context, start, end time.Time) ([]*models.Booking, error) {
	return s.repo.GetByDateRange(ctx, start, end)
}

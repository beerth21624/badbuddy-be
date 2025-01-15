package interfaces

import (
	"context"
	"time"

	"badbuddy/internal/domain/models"

	"github.com/google/uuid"
)

// BookingRepository defines the interface for court booking data operations
type BookingRepository interface {
	Create(ctx context.Context, booking *models.CourtBooking) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.CourtBooking, error)
	List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]models.CourtBooking, error)
	Update(ctx context.Context, booking *models.CourtBooking) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUserBookings(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]models.CourtBooking, error)
	GetVenueBookings(ctx context.Context, venueID uuid.UUID, startDate, endDate time.Time) ([]models.CourtBooking, error)
	GetCourtBookings(ctx context.Context, courtID uuid.UUID, date time.Time) ([]models.CourtBooking, error)
	CheckCourtAvailability(ctx context.Context, courtID uuid.UUID, date time.Time, startTime, endTime time.Time) (bool, error)
	CancelBooking(ctx context.Context, id uuid.UUID) error
	GetPayment(ctx context.Context, bookingID uuid.UUID) (*models.Payment, error)
	CreatePayment(ctx context.Context, payment *models.Payment) error
	UpdatePayment(ctx context.Context, payment *models.Payment) error
	Count(ctx context.Context, userID uuid.UUID, filters map[string]interface{}) (int, error) // Added Count method

}

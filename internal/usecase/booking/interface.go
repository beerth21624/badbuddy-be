package booking

import (
	"context"
	"errors"

	"badbuddy/internal/delivery/dto/requests"
	"badbuddy/internal/delivery/dto/responses"

	"github.com/google/uuid"
)

type UseCase interface {
	CreateBooking(ctx context.Context, userID uuid.UUID, req requests.CreateBookingRequest) (*responses.BookingResponse, error)
	GetBooking(ctx context.Context, id uuid.UUID) (*responses.BookingResponse, error)
	ListBookings(ctx context.Context, userID uuid.UUID, req requests.ListBookingsRequest) (*responses.BookingListResponse, error)
	UpdateBooking(ctx context.Context, id uuid.UUID, req requests.UpdateBookingRequest) (*responses.BookingResponse, error)
	CancelBooking(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetUserBookings(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]responses.BookingResponse, error)
	CheckAvailability(ctx context.Context, req requests.CheckAvailabilityRequest) (*responses.CourtAvailabilityResponse, error)
	GetPayment(ctx context.Context, id uuid.UUID) (*responses.PaymentResponse, error)
	CreatePayment(ctx context.Context, id uuid.UUID, userID uuid.UUID, req requests.CreatePaymentRequest) (*responses.PaymentResponse, error)
	UpdatePayment(ctx context.Context, id uuid.UUID, userID uuid.UUID, req requests.UpdatePaymentRequest) (*responses.PaymentResponse, error)
	ChangeCourtStatus(ctx context.Context) error
}

var (
	ErrUnauthorized = errors.New("unauthorized")

	ErrValidation = errors.New("validation error")

	ErrBookingConflict = errors.New("booking conflict")

	ErrPaymentRequired = errors.New("payment required")

	ErrBookingNotFound = errors.New("booking not found") // Added this line

)

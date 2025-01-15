package models

import (
	"badbuddy/internal/delivery/dto/responses"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BookingStatus string
type PaymentStatus string
type PaymentMethod string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	// BookingStatusCompleted BookingStatus = "completed"

	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"

	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodTransfer PaymentMethod = "transfer"
	PaymentMethodCard     PaymentMethod = "card"
	PaymentMethodQR       PaymentMethod = "qr"
)

// CourtBooking represents a court booking
type CourtBooking struct {
	ID          uuid.UUID     `db:"id"`
	CourtID     uuid.UUID     `db:"court_id"`
	UserID      uuid.UUID     `db:"user_id"`
	Date        time.Time     `db:"booking_date"`
	StartTime   time.Time     `db:"start_time"`
	EndTime     time.Time     `db:"end_time"`
	TotalAmount float64       `db:"total_amount"`
	Status      BookingStatus `db:"status"`
	Notes       *string       `db:"notes"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
	CancelledAt *time.Time    `db:"cancelled_at"`

	// Joined fields
	CourtName     string  `db:"court_name"`
	PricePerHour  float64 `db:"price_per_hour"`
	VenueName     string  `db:"venue_name"`
	VenueLocation string  `db:"venue_location"`
	UserName      string  `db:"user_name"`

	// Related data
	Payment *Payment `db:"-"`
}

// Payment represents a payment for a booking
type Payment struct {
	ID            uuid.UUID     `db:"id"`
	BookingID     uuid.UUID     `db:"booking_id"`
	UserID        uuid.UUID     `db:"user_id"`
	Amount        float64       `db:"amount"`
	Status        PaymentStatus `db:"status"`
	PaymentMethod PaymentMethod `db:"payment_method"`
	TransactionID *string       `db:"transaction_id"`
	CreatedAt     time.Time     `db:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at"`
}

// BookingDetail represents a detailed court booking with all related information
type BookingDetail struct {
	CourtBooking
	Court Court `json:"court"`
	Venue Venue `json:"venue"`
}

// BookingFilters represents the available filters for listing bookings
type BookingFilters struct {
	CourtID  *uuid.UUID     `json:"court_id"`
	VenueID  *uuid.UUID     `json:"venue_id"`
	UserID   *uuid.UUID     `json:"user_id"`
	Date     *time.Time     `json:"date"`
	Status   *BookingStatus `json:"status"`
	DateFrom *time.Time     `json:"date_from"`
	DateTo   *time.Time     `json:"date_to"`
}

// CourtAvailability represents a court's availability for a specific time slot
type CourtAvailability struct {
	CourtID   uuid.UUID `json:"court_id"`
	Date      time.Time `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Available bool      `json:"available"`
}

// Validate validates the booking data
func (b *CourtBooking) Validate() error {
	// Check required fields
	if b.CourtID == uuid.Nil {
		return fmt.Errorf("court ID is required")
	}
	if b.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}
	if b.Date.IsZero() {
		return fmt.Errorf("booking date is required")
	}
	if b.StartTime.IsZero() || b.EndTime.IsZero() {
		return fmt.Errorf("start and end times are required")
	}
	// Check date and time logic
	now := time.Now()
	if b.Date.Before(now.Truncate(24 * time.Hour)) {
		return fmt.Errorf("booking date must be in the future")
	}

	if b.StartTime.After(b.EndTime) {
		return fmt.Errorf("start time must be before end time")
	}

	duration := b.EndTime.Sub(b.StartTime)
	if duration < 30*time.Minute {
		return fmt.Errorf("booking duration must be at least 30 minutes")
	}
	if duration > 4*time.Hour {
		return fmt.Errorf("booking duration cannot exceed 4 hours")
	}

	// Check amount
	if b.TotalAmount <= 0 {
		return fmt.Errorf("total amount must be greater than 0")
	}

	return nil
}

// CalculateDuration returns the booking duration in hours
func (b *CourtBooking) CalculateDuration() float64 {
	duration := b.EndTime.Sub(b.StartTime)
	return duration.Hours()
}

// CalculateTotalAmount calculates the total amount based on duration and price per hour
func (b *CourtBooking) CalculateTotalAmount() float64 {
	duration := b.CalculateDuration()
	return duration * b.PricePerHour
}

// CanBeCancelled checks if the booking can be cancelled based on its status and time
func (b *CourtBooking) CanBeCancelled() bool {
	if b.Status == BookingStatusCancelled || b.Status == BookingStatusConfirmed {
		return false
	}
	return true
	// Check if booking start time is more than 24 hours away
	// now := time.Now()
	// bookingStart := time.Date(
	// 	b.Date.Year(), b.Date.Month(), b.Date.Day(),
	// 	b.StartTime.Hour(), b.StartTime.Minute(), 0, 0, time.Local)
	// return now.Add(24 * time.Hour).Before(bookingStart)
}

// IsOverlapping checks if this booking overlaps with another booking
func (b *CourtBooking) IsOverlapping(other *CourtBooking) bool {
	if b.CourtID != other.CourtID || !b.Date.Equal(other.Date) {
		return false
	}

	return b.StartTime.Before(other.EndTime) && other.StartTime.Before(b.EndTime)
}

// ToResponse converts the booking to a response DTO
func (b *CourtBooking) ToResponse() *responses.BookingResponse {
	resp := &responses.BookingResponse{
		ID:            b.ID.String(),
		CourtName:     b.CourtName,
		VenueName:     b.VenueName,
		VenueLocation: b.VenueLocation,
		UserName:      b.UserName,
		Date:          b.Date.Format("2006-01-02"),
		StartTime:     b.StartTime.Format("15:04"),
		EndTime:       b.EndTime.Format("15:04"),
		TotalAmount:   b.TotalAmount,
		Status:        string(b.Status),
		CreatedAt:     b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     b.UpdatedAt.Format(time.RFC3339),
	}

	if b.Notes != nil {
		resp.Notes = *b.Notes
	}

	if b.CancelledAt != nil {
		resp.CancelledAt = b.CancelledAt.Format(time.RFC3339)
	}

	if b.Payment != nil {
		resp.Payment = &responses.PaymentResponse{
			ID:            b.Payment.ID.String(),
			Amount:        b.Payment.Amount,
			Status:        string(b.Payment.Status),
			PaymentMethod: string(b.Payment.PaymentMethod),
			CreatedAt:     b.Payment.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     b.Payment.UpdatedAt.Format(time.RFC3339),
		}
		if b.Payment.TransactionID != nil {
			resp.Payment.TransactionID = *b.Payment.TransactionID
		}
	}

	return resp
}

// Validate validates the payment data
func (p *Payment) Validate() error {
	if p.BookingID == uuid.Nil {
		return fmt.Errorf("booking ID is required")
	}
	if p.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if p.PaymentMethod == "" {
		return fmt.Errorf("payment method is required")
	}
	return nil
}

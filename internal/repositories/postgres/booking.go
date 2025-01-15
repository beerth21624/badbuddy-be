package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"badbuddy/internal/delivery/dto/responses"
	"badbuddy/internal/domain/models"
	"badbuddy/internal/repositories/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type bookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) interfaces.BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *models.CourtBooking) error {
	// First check availability
	isAvailable, err := r.CheckCourtAvailability(
		ctx,
		booking.CourtID,
		booking.Date,
		booking.StartTime,
		booking.EndTime,
	)
	if err != nil {
		return fmt.Errorf("error checking availability: %w", err)
	}
	if !isAvailable {
		return fmt.Errorf("court is not available for the requested time")
	}

	query := `
        INSERT INTO court_bookings (
            id, court_id, user_id, booking_date, start_time, end_time,
            total_amount, status, notes, created_at, updated_at
        ) VALUES (
            :id, :court_id, :user_id, :booking_date, :start_time, :end_time,
            :total_amount, :status, :notes, :created_at, :updated_at
        )`

	_, err = r.db.NamedExecContext(ctx, query, booking)
	return err
}
func (r *bookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CourtBooking, error) {
	query := `
		SELECT 
			b.*,
			c.name as court_name,
			c.price_per_hour,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as user_name
		FROM court_bookings b
		JOIN courts c ON c.id = b.court_id
		JOIN venues v ON v.id = c.venue_id
		JOIN users u ON u.id = b.user_id
		WHERE b.id = $1`

	var booking models.CourtBooking
	err := r.db.GetContext(ctx, &booking, query, id)
	if err != nil {
		return nil, err
	}

	// Get associated payment if exists
	paymentQuery := `SELECT * FROM payments WHERE booking_id = $1`
	var payment models.Payment
	if err := r.db.GetContext(ctx, &payment, paymentQuery, id); err == nil {
		booking.Payment = &payment
	}

	return &booking, nil
}

func (r *bookingRepository) List(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]models.CourtBooking, error) {
	query := `
		SELECT
			cb.*,
			c.name as court_name,
			c.price_per_hour,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as user_name
		FROM
		users u
		JOIN venues v ON v.owner_id = u.id
		JOIN courts c ON c.venue_id = v.id
		JOIN court_bookings cb ON cb.court_id = c.id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if courtID, ok := filters["court_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND b.court_id = $%d", argCount)
		args = append(args, courtID)
		argCount++
	}

	if date, ok := filters["date"].(string); ok {
		query += fmt.Sprintf(" AND cb.booking_date = $%d", argCount)
		args = append(args, date)
		argCount++
	}

	if status, ok := filters["status"].(string); ok {
		query += fmt.Sprintf(" AND cb.status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if venueID, ok := filters["venue_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND v.id = $%d", argCount)
		args = append(args, venueID)
		argCount++
	}

	if userID != uuid.Nil {
		query += fmt.Sprintf(" AND u.id = $%d", argCount)
		args = append(args, userID)
		argCount++
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
		argCount++
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
		argCount++
	}

	var bookings []models.CourtBooking
	err := r.db.SelectContext(ctx, &bookings, query, args...)
	if err != nil {
		return nil, err
	}

	// Get payments for bookings
	for i, booking := range bookings {
		var payment models.Payment
		paymentQuery := `SELECT * FROM payments WHERE booking_id = $1`
		if err := r.db.GetContext(ctx, &payment, paymentQuery, booking.ID); err == nil {
			bookings[i].Payment = &payment
		}
	}

	return bookings, nil
}

func (r *bookingRepository) Update(ctx context.Context, booking *models.CourtBooking) error {
	query := `
		UPDATE court_bookings SET
			status = :status,
			notes = :notes,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, booking)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

func (r *bookingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM court_bookings WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

func (r *bookingRepository) GetUserBookings(ctx context.Context, userID uuid.UUID, includeHistory bool) ([]models.CourtBooking, error) {
	query := `
		SELECT 
			b.*,
			c.name as court_name,
			c.price_per_hour,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as user_name
		FROM court_bookings b
		JOIN courts c ON c.id = b.court_id
		JOIN venues v ON v.id = c.venue_id
		JOIN users u ON u.id = b.user_id
		WHERE b.user_id = $1`

	if !includeHistory {
		query += " AND b.booking_date >= CURRENT_DATE"
	}

	query += " ORDER BY b.booking_date ASC, b.start_time ASC"

	var bookings []models.CourtBooking
	err := r.db.SelectContext(ctx, &bookings, query, userID)
	if err != nil {
		return nil, err
	}

	// Get payments for bookings
	for i, booking := range bookings {
		var payment models.Payment
		paymentQuery := `SELECT * FROM payments WHERE booking_id = $1`
		if err := r.db.GetContext(ctx, &payment, paymentQuery, booking.ID); err == nil {
			bookings[i].Payment = &payment
		}
	}

	return bookings, nil
}

func (r *bookingRepository) GetVenueBookings(ctx context.Context, venueID uuid.UUID, startDate, endDate time.Time) ([]models.CourtBooking, error) {
	query := `
		SELECT 
			b.*,
			c.name as court_name,
			c.price_per_hour,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as user_name
		FROM court_bookings b
		JOIN courts c ON c.id = b.court_id
		JOIN venues v ON v.id = c.venue_id
		JOIN users u ON u.id = b.user_id
		WHERE v.id = $1 AND b.booking_date BETWEEN $2 AND $3
		ORDER BY b.booking_date ASC, b.start_time ASC`

	var bookings []models.CourtBooking
	err := r.db.SelectContext(ctx, &bookings, query, venueID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get payments for bookings
	for i, booking := range bookings {
		var payment models.Payment
		paymentQuery := `SELECT * FROM payments WHERE booking_id = $1`
		if err := r.db.GetContext(ctx, &payment, paymentQuery, booking.ID); err == nil {
			bookings[i].Payment = &payment
		}
	}

	return bookings, nil
}

func (r *bookingRepository) GetCourtBookings(ctx context.Context, courtID uuid.UUID, date time.Time) ([]models.CourtBooking, error) {
	query := `
		SELECT 
			b.*,
			c.name as court_name,
			c.price_per_hour,
			v.name as venue_name,
			v.location as venue_location,
			u.first_name || ' ' || u.last_name as user_name
		FROM court_bookings b
		JOIN courts c ON c.id = b.court_id
		JOIN venues v ON v.id = c.venue_id
		JOIN users u ON u.id = b.user_id
		WHERE b.court_id = $1 AND b.booking_date = $2
		ORDER BY b.start_time ASC`

	var bookings []models.CourtBooking
	err := r.db.SelectContext(ctx, &bookings, query, courtID, date)
	return bookings, err
}

func (r *bookingRepository) CheckCourtAvailability(ctx context.Context, courtID uuid.UUID, date time.Time, startTime, endTime time.Time) (bool, error) {
	// First check if any existing bookings conflict
	bookingQuery := `
        SELECT COUNT(*)
        FROM court_bookings
        WHERE court_id = $1 
        AND booking_date = $2
        AND status != 'cancelled'
        AND (
            (start_time <= $3 AND end_time > $3)
            OR (start_time < $4 AND end_time >= $4)
            OR (start_time >= $3 AND end_time <= $4)
        )`

	var bookingCount int
	if err := r.db.GetContext(ctx, &bookingCount, bookingQuery, courtID, date, startTime, endTime); err != nil {
		return false, err
	}

	if bookingCount > 0 {
		return false, nil
	}

	// Then check if the venue is open during the requested time
	venueQuery := `
        SELECT v.open_range
        FROM courts c
        JOIN venues v ON v.id = c.venue_id
        WHERE c.id = $1`
	var openRangeJson json.RawMessage
	if err := r.db.GetContext(ctx, &openRangeJson, venueQuery, courtID); err != nil {
		return false, err
	}
	var openRange []responses.OpenRangeResponse
	errParseOpenRangeJson := json.Unmarshal(openRangeJson, &openRange)
	if errParseOpenRangeJson != nil {
		return false, errParseOpenRangeJson
	}
	// Get the day of week for the booking date
	dayOfWeek := strings.ToLower(date.Weekday().String())

	// Check if the requested time falls within venue operating hours
	for _, schedule := range openRange {
		if strings.ToLower(schedule.Day) == dayOfWeek {
			// Convert schedule times to same date as booking for comparison
			scheduleOpen := time.Date(
				startTime.Year(), startTime.Month(), startTime.Day(),
				schedule.OpenTime.Hour(), schedule.OpenTime.Minute(), 0, 0,
				startTime.Location(),
			)
			scheduleClose := time.Date(
				startTime.Year(), startTime.Month(), startTime.Day(),
				schedule.CloseTime.Hour(), schedule.CloseTime.Minute(), 0, 0,
				startTime.Location(),
			)

			// Check if booking time falls within operating hours
			if startTime.Hour() < scheduleOpen.Hour() ||
				(startTime.Hour() == scheduleOpen.Hour() && startTime.Minute() < scheduleOpen.Minute()) ||
				endTime.Hour() > scheduleClose.Hour() ||
				(endTime.Hour() == scheduleClose.Hour() && endTime.Minute() > scheduleClose.Minute()) {
				return false, nil
			}
			return true, nil
		}
	}

	// If we didn't find the day in the schedule, venue is closed
	return false, nil
}

func (r *bookingRepository) CancelBooking(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE court_bookings 
		SET status = 'cancelled', 
			cancelled_at = NOW(), 
			updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

func (r *bookingRepository) GetPayment(ctx context.Context, bookingID uuid.UUID) (*models.Payment, error) {
	query := `SELECT * FROM payments WHERE booking_id = $1`

	var payment models.Payment
	err := r.db.GetContext(ctx, &payment, query, bookingID)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *bookingRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (
			id, booking_id, user_id, amount, status, payment_method,
			transaction_id, created_at, updated_at
		) VALUES (
			:id, :booking_id, :user_id, :amount, :status, :payment_method,
			:transaction_id, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, payment)
	return err
}

func (r *bookingRepository) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	query := `
		UPDATE payments SET
			status = :status,
			payment_method = :payment_method,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, payment)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

func (r *bookingRepository) Count(ctx context.Context, userID uuid.UUID, filters map[string]interface{}) (int, error) {
	query := `
		SELECT
			COUNT(*)
		FROM
		users u
		JOIN venues v ON v.owner_id = u.id
		JOIN courts c ON c.venue_id = v.id
		JOIN court_bookings cb ON cb.court_id = c.id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if courtID, ok := filters["court_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND b.court_id = $%d", argCount)
		args = append(args, courtID)
		argCount++
	}

	if date, ok := filters["date"].(time.Time); ok {
		query += fmt.Sprintf(" AND b.booking_date = $%d", argCount)
		args = append(args, date)
		argCount++
	}

	if status, ok := filters["status"].(string); ok {
		query += fmt.Sprintf(" AND b.status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	if venueID, ok := filters["venue_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND v.id = $%d", argCount)
		args = append(args, venueID)
		argCount++
	}

	if userID != uuid.Nil {
		query += fmt.Sprintf(" AND u.id = $%d", argCount)
		args = append(args, userID)
		argCount++
	}
	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

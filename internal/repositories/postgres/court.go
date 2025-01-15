package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"badbuddy/internal/domain/models"
	"badbuddy/internal/repositories/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type courtRepository struct {
	db *sqlx.DB
}

func NewCourtRepository(db *sqlx.DB) interfaces.CourtRepository {
	return &courtRepository{db: db}
}

func (r *courtRepository) Create(ctx context.Context, court *models.Court) error {
	query := `
		INSERT INTO courts (
			id, venue_id, name, description, price_per_hour,
			status, created_at, updated_at
		) VALUES (
			:id, :venue_id, :name, :description, :price_per_hour,
			:status, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, court)
	return err
}

func (r *courtRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Court, error) {
	query := `
		SELECT 
			*
		FROM courts
		WHERE id = $1 AND deleted_at IS NULL`

	var court models.Court
	err := r.db.GetContext(ctx, &court, query, id)
	if err != nil {
		return nil, err
	}

	return &court, nil
}

func (r *courtRepository) GetCourtWithVenueByID(ctx context.Context, id uuid.UUID) (*models.CourtWithVenue, error) {
	query := `
		SELECT 
			c.*,
			v.name as venue_name,
			v.location as venue_location,
			v.status as venue_status
		FROM courts c
		JOIN venues v ON v.id = c.venue_id
		WHERE c.id = $1 AND c.deleted_at IS NULL`

	var court models.CourtWithVenue
	err := r.db.GetContext(ctx, &court, query, id)
	if err != nil {
		return nil, err
	}

	return &court, nil
}

func (r *courtRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Court, error) {
	query := `
		SELECT 
			c.*,
			v.name as venue_name,
			v.location as venue_location,
			v.status as venue_status
		FROM courts c
		JOIN venues v ON v.id = c.venue_id
		WHERE c.deleted_at IS NULL`

	args := []interface{}{}
	argCount := 1

	// Add filters
	if len(filters) > 0 {
		whereConditions := []string{}

		if venueID, ok := filters["venue_id"].(uuid.UUID); ok {
			whereConditions = append(whereConditions, fmt.Sprintf("c.venue_id = $%d", argCount))
			args = append(args, venueID)
			argCount++
		}

		if status, ok := filters["status"].(models.CourtStatus); ok {
			whereConditions = append(whereConditions, fmt.Sprintf("c.status = $%d", argCount))
			args = append(args, status)
			argCount++
		}

		if location, ok := filters["location"].(string); ok {
			whereConditions = append(whereConditions, fmt.Sprintf("v.location ILIKE $%d", argCount))
			args = append(args, "%"+location+"%")
			argCount++
		}

		if priceMin, ok := filters["price_min"].(float64); ok {
			whereConditions = append(whereConditions, fmt.Sprintf("c.price_per_hour >= $%d", argCount))
			args = append(args, priceMin)
			argCount++
		}

		if priceMax, ok := filters["price_max"].(float64); ok {
			whereConditions = append(whereConditions, fmt.Sprintf("c.price_per_hour <= $%d", argCount))
			args = append(args, priceMax)
			argCount++
		}

		if len(whereConditions) > 0 {
			query += " AND " + strings.Join(whereConditions, " AND ")
		}
	}

	// Add ordering
	query += " ORDER BY c.created_at DESC"

	// Add pagination
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
	var courts []models.Court
	err := r.db.SelectContext(ctx, &courts, query, args...)
	if err != nil {
		return nil, err
	}

	return courts, nil
}

func (r *courtRepository) Update(ctx context.Context, court *models.Court) error {
	query := `
		UPDATE courts SET
			name = :name,
			description = :description,
			price_per_hour = :price_per_hour,
			status = :status,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`

	result, err := r.db.NamedExecContext(ctx, query, court)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("court not found")
	}

	return nil
}

func (r *courtRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE courts SET
			deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("court not found")
	}

	return nil
}

func (r *courtRepository) GetByVenue(ctx context.Context, venueID uuid.UUID) ([]models.Court, error) {
	query := `
		SELECT 
			*
		FROM courts
		WHERE venue_id = $1 AND deleted_at IS NULL
		ORDER BY c.name ASC`

	var courts []models.Court
	err := r.db.SelectContext(ctx, &courts, query, venueID)
	return courts, err
}

func (r *courtRepository) GetCourtWithVenueByVenue(ctx context.Context, venueID uuid.UUID) ([]models.CourtWithVenue, error) {
	query := `
		SELECT 
			c.*,
			v.name as venue_name,
			v.location as venue_location,
			v.status as venue_status
		FROM courts c
		JOIN venues v ON v.id = c.venue_id
		WHERE c.venue_id = $1 AND c.deleted_at IS NULL
		ORDER BY c.name ASC`

	var courts []models.CourtWithVenue
	err := r.db.SelectContext(ctx, &courts, query, venueID)
	return courts, err
}

func (r *courtRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.CourtStatus) error {
	query := `
		UPDATE courts SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("court not found")
	}

	return nil
}

func (r *courtRepository) GetAvailableCourts(ctx context.Context, venueID uuid.UUID, date time.Time, startTime, endTime time.Time) ([]models.Court, error) {
	query := `
		SELECT DISTINCT
			c.*,
			v.name as venue_name,
			v.location as venue_location,
			v.status as venue_status
		FROM courts c
		JOIN venues v ON v.id = c.venue_id
		WHERE c.venue_id = $1 
		AND c.deleted_at IS NULL
		AND c.status = 'available'
		AND v.status = 'active'
		AND NOT EXISTS (
			SELECT 1 FROM court_bookings b
			WHERE b.court_id = c.id
			AND b.booking_date = $2
			AND b.status != 'cancelled'
			AND (
				(b.start_time <= $3 AND b.end_time > $3)
				OR (b.start_time < $4 AND b.end_time >= $4)
				OR (b.start_time >= $3 AND b.end_time <= $4)
			)
		)
		ORDER BY c.name ASC`

	var courts []models.Court
	err := r.db.SelectContext(ctx, &courts, query, venueID, date, startTime, endTime)
	return courts, err
}

func (r *courtRepository) Count(ctx context.Context, filters map[string]interface{}) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM courts c
		JOIN venues v ON v.id = c.venue_id
		WHERE c.deleted_at IS NULL`

	args := []interface{}{}
	argCount := 1

	// Add filters
	if len(filters) > 0 {
		whereConditions := []string{}

		if venueID, ok := filters["venue_id"].(uuid.UUID); ok {
			whereConditions = append(whereConditions, fmt.Sprintf("c.venue_id = $%d", argCount))
			args = append(args, venueID)
			argCount++
		}

		if status, ok := filters["status"].(models.CourtStatus); ok {
			whereConditions = append(whereConditions, fmt.Sprintf("c.status = $%d", argCount))
			args = append(args, status)
			argCount++
		}

		if location, ok := filters["location"].(string); ok {
			whereConditions = append(whereConditions, fmt.Sprintf("v.location ILIKE $%d", argCount))
			args = append(args, "%"+location+"%")
			argCount++
		}

		if len(whereConditions) > 0 {
			query += " AND " + strings.Join(whereConditions, " AND ")
		}
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

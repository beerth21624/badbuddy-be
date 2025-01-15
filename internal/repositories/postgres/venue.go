package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"badbuddy/internal/domain/models"
	"badbuddy/internal/repositories/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type venueRepository struct {
	db *sqlx.DB
}

func NewVenueRepository(db *sqlx.DB) interfaces.VenueRepository {
	return &venueRepository{
		db: db,
	}
}

func (r *venueRepository) Create(ctx context.Context, venue *models.Venue) error {
	// First check if venue with same name exists
	checkQuery := `
        SELECT EXISTS (
            SELECT 1 FROM venues 
            WHERE LOWER(name) = LOWER($1)
        )
    `

	var exists bool
	err := r.db.GetContext(ctx, &exists, checkQuery, venue.Name)
	if err != nil {
		return fmt.Errorf("failed to check venue name: %w", err)
	}

	if exists {
		return fmt.Errorf("venue with name '%s' already exists", venue.Name)
	}

	venueInsert := models.VenueInsert{
		ID:            venue.ID,
		Name:          venue.Name,
		Description:   venue.Description,
		Address:       venue.Address,
		Location:      venue.Location,
		Phone:         venue.Phone,
		Email:         venue.Email,
		OpenRange:     venue.OpenRange.RawMessage,
		ImageURLs:     venue.ImageURLs,
		Status:        venue.Status,
		Rating:        venue.Rating,
		TotalReviews:  venue.TotalReviews,
		OwnerID:       venue.OwnerID,
		CreatedAt:     venue.CreatedAt,
		UpdatedAt:     venue.UpdatedAt,
		DeletedAt:     venue.DeletedAt,
		Search_vector: venue.Search_vector,
		Rules:         venue.Rules.RawMessage,
		Facilities:    venue.Facilities,
		Latitude:      venue.Latitude,
		Longitude:     venue.Longitude,
	}

	// If no duplicate, proceed with insert
	insertQuery := `
        INSERT INTO venues (
            id, name, description, address, location, phone, email,
            open_range, image_urls, status, rating,
            total_reviews, owner_id, created_at, updated_at, rules, latitude, longitude
        ) VALUES (
            safe_generate_uuid(), :name, :description, :address, :location, :phone, :email,
            :open_range, :image_urls, :status, :rating,
            :total_reviews, :owner_id, :created_at, :updated_at, :rules, :latitude, :longitude
        )
        RETURNING *
    `

	// Use NamedQueryRow instead of NamedExec to get the returned values
	rows, err := r.db.NamedQueryContext(ctx, insertQuery, venueInsert)
	if err != nil {
		return fmt.Errorf("failed to create venue: %w", err)
	}
	defer rows.Close()

	// Scan the returned row into the venue struct
	if rows.Next() {
		err = rows.StructScan(venue)
		if err != nil {
			return fmt.Errorf("failed to scan venue: %w", err)
		}
	}

	return nil
}

// You might want to create custom errors for better error handling
type ErrDuplicateVenue struct {
	Name string
}

func (e *ErrDuplicateVenue) Error() string {
	return fmt.Sprintf("venue with name '%s' already exists", e.Name)
}

// Optional: Create a method to check if an error is a duplicate venue error
func IsDuplicateVenueError(err error) bool {
	_, ok := err.(*ErrDuplicateVenue)
	return ok
}

func (r *venueRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.VenueWithCourts, error) {
	result := &models.VenueWithCourts{}

	// Get venue details
	query := `
		SELECT * FROM venues  WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.GetContext(ctx, &result.Venue, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("venue not found")
		}
		return nil, fmt.Errorf("failed to get venue: %w", err)
	}

	// Get facilities
	facilitiesQuery := `
		SELECT f.id, f.name
		FROM venues_facilities vf
		JOIN facilities f ON vf.facility_id = f.id
		WHERE venue_id = $1`

	err = r.db.SelectContext(ctx, &result.Venue.Facilities, facilitiesQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get facilities: %w", err)
	}

	courtsQuery := `
		SELECT * FROM courts 
		WHERE venue_id = $1 AND deleted_at IS NULL 
		ORDER BY created_at`
	err = r.db.SelectContext(ctx, &result.Courts, courtsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get courts: %w", err)
	}

	return result, nil
}

func (r *venueRepository) Update(ctx context.Context, venue *models.Venue) error {

	params := map[string]interface{}{
		"id":          venue.ID,
		"name":        venue.Name,
		"description": venue.Description,
		"address":     venue.Address,
		"location":    venue.Location,
		"phone":       venue.Phone,
		"email":       venue.Email,
		"open_range":  venue.OpenRange.RawMessage,
		"image_urls":  venue.ImageURLs,
		"status":      venue.Status,
		"updated_at":  venue.UpdatedAt,
		"rules":       venue.Rules.RawMessage,
		"latitude":    venue.Latitude,
		"longitude":   venue.Longitude,
	}

	query := `
		UPDATE venues SET
			name = :name,
			description = :description,
			address = :address,
			location = :location,
			phone = :phone,
			email = :email,
			open_range = :open_range,
			image_urls = :image_urls,
			status = :status,
			updated_at = :updated_at,
			rules = :rules,
			latitude = :latitude,
			longitude = :longitude
		WHERE id = :id AND deleted_at IS NULL`

	result, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update venue: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("venue not found")
	}

	return nil
}

func (r *venueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE venues 
		SET deleted_at = NOW(), updated_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete venue: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("venue not found")
	}

	return nil
}

func (r *venueRepository) List(ctx context.Context, location string, limit, offset int) ([]models.Venue, error) {
	query := `
		SELECT 
			v.id, v.name, v.description, v.address, v.location, v.phone, v.email,
			v.open_range, v.image_urls, v.status, v.rating, v.total_reviews, v.owner_id,
			v.created_at, v.updated_at, v.search_vector, v.rules,
			COALESCE(json_agg(
				json_build_object('id', f.id, 'name', f.name)
			) FILTER (WHERE f.id IS NOT NULL), '[]') AS facilities,
			COALESCE(json_agg(
				json_build_object('id', c.id, 'name', c.name, 'description', c.description, 'price_per_hour', c.price_per_hour, 'status', c.status)
			) FILTER (WHERE c.id IS NOT NULL), '[]') AS courts
		FROM 
			venues v
		LEFT JOIN 
			venues_facilities vf ON v.id = vf.venue_id
		LEFT JOIN 
			facilities f ON vf.facility_id = f.id
		LEFT JOIN
			courts c ON v.id = c.venue_id
		WHERE 
			v.deleted_at IS NULL
			AND ($1 = '' OR v.location = $1)
		GROUP BY 
			v.id
		ORDER BY 
			v.rating DESC, v.total_reviews DESC, v.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, location, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list venues: %w", err)
	}
	defer rows.Close()

	var venues []models.Venue
	for rows.Next() {
		var venue models.Venue
		var facilitiesJSON []byte
		var courtsJSON []byte

		err := rows.Scan(
			&venue.ID, &venue.Name, &venue.Description, &venue.Address, &venue.Location,
			&venue.Phone, &venue.Email, &venue.OpenRange, &venue.ImageURLs,
			&venue.Status, &venue.Rating, &venue.TotalReviews, &venue.OwnerID,
			&venue.CreatedAt, &venue.UpdatedAt, &venue.Search_vector, &venue.Rules,
			&facilitiesJSON, &courtsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan venue: %w", err)
		}

		// Unmarshal facilities JSON into the Facilities slice
		err = json.Unmarshal(facilitiesJSON, &venue.Facilities)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal facilities for venue %s: %w", venue.ID, err)
		}

		// Unmarshal courts JSON into the Courts slice
		err = json.Unmarshal(courtsJSON, &venue.Courts)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal courts for venue %s: %w", venue.ID, err)
		}

		venues = append(venues, venue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return venues, nil
}

func (r *venueRepository) CountVenues(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM venues 
		WHERE deleted_at IS NULL
`
	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count venues: %w", err)
	}

	return count, nil
}

func (r *venueRepository) Search(ctx context.Context, query string, limit, offset int, minPrice, maxPrice int, location string, facilities []string) ([]models.Venue, error) {
	searchQuery := `
		SELECT 
			v.id, v.name, v.description, v.address, v.location, v.phone, v.email,
			v.open_range, v.image_urls, v.status, v.rating, v.total_reviews, v.owner_id,
			v.created_at, v.updated_at, v.rules, v.latitude, v.longitude,
			COALESCE(json_agg(
				json_build_object('id', f.id, 'name', f.name)
			) FILTER (WHERE f.id IS NOT NULL), '[]') AS facilities,
			COALESCE(json_agg(
				json_build_object('id', c.id, 'name', c.name, 'description', c.description, 'price_per_hour', c.price_per_hour, 'status', c.status)
			) FILTER (WHERE c.id IS NOT NULL), '[]') AS courts
		FROM 
			venues v
		LEFT JOIN 
			venues_facilities vf ON v.id = vf.venue_id
		LEFT JOIN 
			facilities f ON vf.facility_id = f.id
		LEFT JOIN
			courts c ON v.id = c.venue_id
		WHERE 
			v.deleted_at IS NULL
			AND (
				v.search_vector @@ plainto_tsquery($1)
				OR v.name ILIKE '%' || $1 || '%'
			)
			AND c.price_per_hour >= $3
			AND c.price_per_hour <= $4
			AND ($2 = '' OR v.location = $2)`

	// Add facilities filter if provided
	if len(facilities) > 0 {
		placeholders := make([]string, len(facilities))
		for i := range facilities {
			placeholders[i] = fmt.Sprintf("$%d", i+7)
		}
		facilitiesCondition := fmt.Sprintf(
			"AND v.id IN (SELECT venue_id FROM venues_facilities vf2 JOIN facilities f2 ON vf2.facility_id = f2.id WHERE f2.name IN (%s) GROUP BY venue_id HAVING COUNT(DISTINCT f2.name) = %d)",
			strings.Join(placeholders, ", "), len(facilities),
		)
		searchQuery += " " + facilitiesCondition
	}

	// Close the query with GROUP BY and ORDER BY clauses
	searchQuery += `
		GROUP BY 
			v.id
		ORDER BY 
			v.rating DESC, v.total_reviews DESC, v.created_at DESC
		LIMIT $5 OFFSET $6`

	// Prepare parameters, including facilities
	params := []interface{}{query, location, minPrice, maxPrice, limit, offset}
	for _, facility := range facilities {
		params = append(params, facility)
	}

	// Execute the query
	rows, err := r.db.QueryContext(ctx, searchQuery, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to search venues: %w", err)
	}
	defer rows.Close()

	var venues []models.Venue
	for rows.Next() {
		var venue models.Venue
		var facilitiesJSON []byte
		var courtsJSON []byte

		// Scan venue fields, then the aggregated JSON for facilities
		err := rows.Scan(
			&venue.ID, &venue.Name, &venue.Description, &venue.Address, &venue.Location,
			&venue.Phone, &venue.Email, &venue.OpenRange, &venue.ImageURLs,
			&venue.Status, &venue.Rating, &venue.TotalReviews, &venue.OwnerID,
			&venue.CreatedAt, &venue.UpdatedAt, &venue.Rules, &venue.Latitude, &venue.Longitude,
			&facilitiesJSON, &courtsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan venue: %w", err)
		}

		// Unmarshal facilities JSON into the Facilities slice
		err = json.Unmarshal(facilitiesJSON, &venue.Facilities)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal facilities for venue %s: %w", venue.ID, err)
		}

		// Unmarshal courts JSON into the Courts slice
		err = json.Unmarshal(courtsJSON, &venue.Courts)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal courts for venue %s: %w", venue.ID, err)
		}

		venues = append(venues, venue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return venues, nil
}

func (r *venueRepository) CountSearch(ctx context.Context, query string, minPrice, maxPrice int, location string, facilities []string) (int, error) {
	countQuery := `
		SELECT 
			COUNT(DISTINCT v.id)
		FROM 
			venues v
		LEFT JOIN 
			venues_facilities vf ON v.id = vf.venue_id
		LEFT JOIN 
			facilities f ON vf.facility_id = f.id
		LEFT JOIN
			courts c ON v.id = c.venue_id
		WHERE 
			v.deleted_at IS NULL
			AND (
				v.search_vector @@ plainto_tsquery($1)
				OR v.name ILIKE '%' || $1 || '%'
			)
			AND c.price_per_hour >= $3
			AND c.price_per_hour <= $4
			AND ($2 = '' OR v.location = $2)`

	// Add facilities filter if provided
	if len(facilities) > 0 {
		placeholders := make([]string, len(facilities))
		for i := range facilities {
			placeholders[i] = fmt.Sprintf("$%d", i+5)
		}
		facilitiesCondition := fmt.Sprintf(
			"AND v.id IN (SELECT venue_id FROM venues_facilities vf2 JOIN facilities f2 ON vf2.facility_id = f2.id WHERE f2.name IN (%s) GROUP BY venue_id HAVING COUNT(DISTINCT f2.name) = %d)",
			strings.Join(placeholders, ", "), len(facilities),
		)
		countQuery += " " + facilitiesCondition
	}

	// Prepare parameters, including facilities
	params := []interface{}{query, location, minPrice, maxPrice}
	for _, facility := range facilities {
		params = append(params, facility)
	}

	// Execute the count query
	var count int
	err := r.db.QueryRowContext(ctx, countQuery, params...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count venues: %w", err)
	}

	return count, nil
}


func (r *venueRepository) AddCourt(ctx context.Context, court *models.Court) error {
	query := `
		INSERT INTO courts (
			id, venue_id, name, description, price_per_hour,
			status, created_at, updated_at
		) VALUES (
			:id, :venue_id, :name, :description, :price_per_hour,
			:status, :created_at, :updated_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, court)
	if err != nil {
		return fmt.Errorf("failed to add court: %w", err)
	}

	return nil
}

func (r *venueRepository) UpdateCourt(ctx context.Context, court *models.Court) error {
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
		return fmt.Errorf("failed to update court: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("court not found")
	}

	return nil
}

func (r *venueRepository) DeleteCourt(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE courts 
		SET deleted_at = NOW(), updated_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete court: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("court not found")
	}

	return nil
}

func (r *venueRepository) GetCourts(ctx context.Context, venueID uuid.UUID) ([]models.Court, error) {
	query := `
		SELECT * FROM courts 
		WHERE venue_id = $1 AND deleted_at IS NULL 
		ORDER BY created_at`

	courts := []models.Court{}
	err := r.db.SelectContext(ctx, &courts, query, venueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get courts: %w", err)
	}

	return courts, nil
}

func (r *venueRepository) AddReview(ctx context.Context, review *models.VenueReview) error {

	// Insert review
	query := `
		INSERT INTO venue_reviews (
			id, venue_id, user_id, rating, comment, created_at
		) VALUES (
			:id, :venue_id, :user_id, :rating, :comment, :created_at
		)`

	_, err := r.db.NamedExecContext(ctx, query, review)
	if err != nil {
		return fmt.Errorf("failed to add review: %w", err)
	}

	fmt.Println(review)

	// Update venue rating
	err = r.UpdateVenueRating(ctx, review.VenueID)
	if err != nil {
		return fmt.Errorf("failed to update venue rating: %w", err)
	}

	return nil
}

func (r *venueRepository) GetReviews(ctx context.Context, venueID uuid.UUID, limit, offset int) ([]models.VenueReview, error) {
	query := `
		SELECT vr.*, 
			u.id as user_id
		FROM venue_reviews vr
		JOIN users u ON u.id = vr.user_id
		WHERE vr.venue_id = $1
		ORDER BY vr.created_at DESC
		LIMIT $2 OFFSET $3`

	reviews := []models.VenueReview{}
	err := r.db.SelectContext(ctx, &reviews, query, venueID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, nil
}

func (r *venueRepository) UpdateVenueRating(ctx context.Context, venueID uuid.UUID) error {
	query := `
		UPDATE venues 
		SET 
			rating = (
				SELECT COALESCE(AVG(rating)::NUMERIC(3,2), 0)
				FROM venue_reviews
				WHERE venue_id = $1
			),
			total_reviews = (
				SELECT COUNT(*)
				FROM venue_reviews
				WHERE venue_id = $1
			),
			updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, venueID)
	if err != nil {
		return fmt.Errorf("failed to update venue rating: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("venue not found")
	}

	return nil
}

func (r *venueRepository) GetFacilities(ctx context.Context, venueID uuid.UUID) ([]models.Facility, error) {
	query := `
		SELECT
			f.id,
			f.name
		FROM
			venues_facilities vf
			JOIN facilities f ON vf.facility_id = f.id
		WHERE
			venue_id = $1;`

	facilities := []models.Facility{}
	err := r.db.SelectContext(ctx, &facilities, query, venueID)

	if err != nil {
		return nil, fmt.Errorf("failed to get facilities: %w", err)
	}

	return facilities, nil
}

func (r *venueRepository) AddFacilities(ctx context.Context, venueID uuid.UUID, facilityIDs []uuid.UUID) error {
	query := `
		INSERT INTO venues_facilities (venue_id, facility_id)
		VALUES (:venue_id, :facility_id)`

	for i := range facilityIDs {
		_, err := r.db.NamedExecContext(ctx, query, map[string]interface{}{
			"venue_id":    venueID,
			"facility_id": facilityIDs[i],
		})
		if err != nil {
			return fmt.Errorf("failed to add facility: %w", err)
		}
	}

	return nil
}

func (r *venueRepository) UpdateFacilities(ctx context.Context, venueID uuid.UUID, facilityIDs []uuid.UUID) error {
	// Delete all existing facilities
	deleteQuery := `
		DELETE FROM venues_facilities
		WHERE venue_id = $1`
	_, err := r.db.ExecContext(ctx, deleteQuery, venueID)
	if err != nil {
		return fmt.Errorf("failed to update facilities: %w", err)
	}

	err = r.AddFacilities(ctx, venueID, facilityIDs)
	if err != nil {
		return fmt.Errorf("failed to update facilities: %w", err)
	}

	return nil
}

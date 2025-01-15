package postgres

import (
	"badbuddy/internal/domain/models"
	"badbuddy/internal/repositories/interfaces"
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type facilityRepository struct {
	db *sqlx.DB
}

func NewFacilityRepository(db *sqlx.DB) interfaces.FacilityRepository {
	return &facilityRepository{db: db}
}

func (r *facilityRepository) GetFacilities(ctx context.Context) ([]models.Facility, error) {
	facilities := []models.Facility{}

	query := `SELECT * FROM facilities`

	err := r.db.SelectContext(ctx, &facilities, query)
	if err != nil {
		return nil, err
	}

	return facilities, nil
}

func (r *facilityRepository) GetFacilityByID(ctx context.Context, id uuid.UUID) (*models.Facility, error) {
	facility := models.Facility{}

	query := `SELECT * FROM facilities WHERE id = $1`

	err := r.db.GetContext(ctx, &facility, query, id)
	if err != nil {
		return nil, err
	}

	return &facility, nil
}

func (r *facilityRepository) CreateFacility(ctx context.Context, facility *models.Facility) error {
	query := `INSERT INTO facilities (id, name) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, facility.ID, facility.Name)
	if err != nil {
		return err
	}

	return nil
}

func (r *facilityRepository) UpdateFacility(ctx context.Context, facility *models.Facility) error {
	query := `UPDATE facilities SET name = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, facility.Name, facility.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *facilityRepository) DeleteFacility(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM facilities WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
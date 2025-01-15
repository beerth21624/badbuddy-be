// internal/domain/models/venue.go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type VenueStatus string
type CourtStatus string

const (
	VenueStatusActive      VenueStatus = "active"
	VenueStatusInactive    VenueStatus = "inactive"
	VenueStatusMaintenance VenueStatus = "maintenance"

	CourtStatusAvailable   CourtStatus = "available"
	CourtStatusOccupied    CourtStatus = "occupied"
	CourtStatusMaintenance CourtStatus = "maintenance"
)

// NullRawMessage is a custom type that properly handles NULL JSON values
type NullRawMessage struct {
	json.RawMessage
	Valid bool
}

// Value implements the driver.Valuer interface
func (n NullRawMessage) Value() (driver.Value, error) {
	if !n.Valid || len(n.RawMessage) == 0 {
		return nil, nil
	}
	return n.RawMessage, nil
}

// Scan implements the sql.Scanner interface
func (n *NullRawMessage) Scan(value interface{}) error {
	if value == nil {
		n.RawMessage, n.Valid = nil, false
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if len(v) == 0 {
			n.RawMessage, n.Valid = nil, false
			return nil
		}
		n.RawMessage = json.RawMessage(v)
		n.Valid = true
		return nil
	case string:
		if v == "" {
			n.RawMessage, n.Valid = nil, false
			return nil
		}
		n.RawMessage = json.RawMessage(v)
		n.Valid = true
		return nil
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type *NullRawMessage", value)
	}
}

// MarshalJSON implements json.Marshaler
func (n NullRawMessage) MarshalJSON() ([]byte, error) {
	if !n.Valid || len(n.RawMessage) == 0 {
		return []byte("null"), nil
	}
	return n.RawMessage, nil
}

// UnmarshalJSON implements json.Unmarshaler
func (n *NullRawMessage) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		n.Valid = false
		n.RawMessage = nil
		return nil
	}
	n.Valid = true
	n.RawMessage = json.RawMessage(data)
	return nil
}

type Venue struct {
	ID            uuid.UUID      `db:"id"`
	Name          string         `db:"name"`
	Description   string         `db:"description"`
	Address       string         `db:"address"`
	Location      string         `db:"location"`
	Phone         string         `db:"phone"`
	Email         string         `db:"email"`
	OpenRange     NullRawMessage `db:"open_range"`
	ImageURLs     string         `db:"image_urls"`
	Status        VenueStatus    `db:"status"`
	Rating        float64        `db:"rating"`
	TotalReviews  int            `db:"total_reviews"`
	OwnerID       uuid.UUID      `db:"owner_id"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
	DeletedAt     *time.Time     `db:"deleted_at"`
	Search_vector string         `db:"search_vector"`
	Rules         NullRawMessage `db:"rules"`
	Facilities    []Facility     `db:"facilities"`
	Courts        []Court        `db:"courts"`
	Latitude      float64        `db:"latitude"`
	Longitude     float64        `db:"longitude"`
}
type VenueInsert struct {
	ID            uuid.UUID   `db:"id"`
	Name          string      `db:"name"`
	Description   string      `db:"description"`
	Address       string      `db:"address"`
	Location      string      `db:"location"`
	Phone         string      `db:"phone"`
	Email         string      `db:"email"`
	OpenRange     []byte      `db:"open_range"`
	ImageURLs     string      `db:"image_urls"`
	Status        VenueStatus `db:"status"`
	Rating        float64     `db:"rating"`
	TotalReviews  int         `db:"total_reviews"`
	OwnerID       uuid.UUID   `db:"owner_id"`
	CreatedAt     time.Time   `db:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at"`
	DeletedAt     *time.Time  `db:"deleted_at"`
	Search_vector string      `db:"search_vector"`
	Rules         []byte      `db:"rules"`
	Facilities    []Facility  `db:"facilities"`
	Latitude      float64     `db:"latitude"`
	Longitude     float64     `db:"longitude"`
}

type Court struct {
	ID            uuid.UUID   `db:"id"`
	VenueID       uuid.UUID   `db:"venue_id"`
	VenueName     string      `db:"venue_name"`
	VenueLocation string      `db:"venue_location"`
	VenueStatus   VenueStatus `db:"venue_status"`
	Name          string      `db:"name"`
	Description   string      `db:"description"`
	PricePerHour  float64     `db:"price_per_hour"`
	Status        CourtStatus `db:"status"`
	CreatedAt     time.Time   `db:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at"`
	DeletedAt     *time.Time  `db:"deleted_at"`
}

type VenueWithCourts struct {
	Venue
	Courts []Court `db:"courts"`
}

type VenueReview struct {
	ID        uuid.UUID `db:"id"`
	VenueID   uuid.UUID `db:"venue_id"`
	UserID    uuid.UUID `db:"user_id"`
	Rating    int       `db:"rating"`
	Comment   string    `db:"comment"`
	CreatedAt time.Time `db:"created_at"`
	UpdateAt  time.Time `db:"updated_at"`
}

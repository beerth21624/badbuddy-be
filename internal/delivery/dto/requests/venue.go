package requests

import (
	"time"
)

type CreateVenueRequest struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description"`
	Address     string      `json:"address" validate:"required"`
	Location    string      `json:"location" validate:"required"`
	Phone       string      `json:"phone" validate:"required"`
	Email       string      `json:"email" validate:"required,email"`
	OpenRange   []OpenRange `json:"open_range" validate:"required"`
	ImageURLs   string      `json:"image_urls"`
	Rules       []Rule      `json:"rules"`
	Facilities  []Facility  `json:"facilities" validate:"required"`
	Latitude    float64     `json:"latitude"`
	Longitude   float64     `json:"longitude"`
}

type Facility struct {
	ID string `json:"id"`
}

type OpenRange struct {
	Day       string    `json:"day"`
	IsOpen    bool      `json:"is_open"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
}

type Rule struct {
	Rule string `json:"rule"`
}

type UpdateVenueRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Address     string      `json:"address"`
	Location    string      `json:"location"`
	Phone       string      `json:"phone"`
	Email       string      `json:"email"`
	OpenRange   []OpenRange `json:"open_range" validate:"required"`
	ImageURLs   string      `json:"image_urls"`
	Status      string      `json:"status"`
	Rules       []Rule      `json:"rules"`
	Facilities  []Facility  `json:"facilities"`
	Latitude    float64     `json:"latitude"`
	Longitude   float64     `json:"longitude"`
}

// type CreateCourtRequest struct {
// 	Name         string  `json:"name" validate:"required"`
// 	Description  string  `json:"description"`
// 	PricePerHour float64 `json:"price_per_hour" validate:"required,gt=0"`
// }

// type UpdateCourtRequest struct {
// 	CourtID      string  `json:"court_id" validate:"required,uuid"`
// 	Name         string  `json:"name"`
// 	Description  string  `json:"description"`
// 	PricePerHour float64 `json:"price_per_hour" validate:"gt=0"`
// 	Status       string  `json:"status"`
// }

type AddReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment"`
}
